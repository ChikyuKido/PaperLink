package pvf

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"paperlink/util"
)

var log = util.GroupLog("PVF")

func WritePVFFromPDF(inputFile string) (string, error) {
	start := time.Now()
	log.Infof("WritePVFFromPDF start input=%s", inputFile)

	tempDir, err := os.MkdirTemp(os.TempDir(), "pvf_*")
	if err != nil {
		return "", fmt.Errorf("could not create temporary directory: %w", err)
	}

	splitStart := time.Now()
	if err := api.SplitFile(inputFile, tempDir, 1, nil); err != nil {
		return "", fmt.Errorf("pdfcpu failed to split file: %w", err)
	}
	log.Infof("WritePVFFromPDF split done dir=%s took=%s", tempDir, time.Since(splitStart))

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return "", fmt.Errorf("could not read temporary directory: %w", err)
	}
	sort.Slice(files, func(i, j int) bool {
		return getNum(files[i].Name()) < getNum(files[j].Name())
	})

	pageCount := countEntriesByExt(files, ".pdf")
	outputFilePath := fmt.Sprintf("%s/output.pvf", tempDir)
	writeStart := time.Now()
	if err := writePVFByFileEntries(tempDir, files, ".pdf", outputFilePath); err != nil {
		return "", err
	}
	log.Infof("WritePVFFromPDF done pages=%d out=%s write=%s total=%s", pageCount, outputFilePath, time.Since(writeStart), time.Since(start))
	return outputFilePath, nil
}

func writePVFByFileEntries(tempDir string, files []os.DirEntry, fileExt string, outputFilePath string) error {
	start := time.Now()
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer outputFile.Close()

	pageCount := uint64(0)
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != fileExt {
			continue
		}
		pageCount++
	}

	outputFile.Write(PVF_MAGIC_BYTES)
	outputFile.Write(VERSION)
	binary.Write(outputFile, binary.LittleEndian, pageCount)

	currentOffset := uint64(13)
	offsetForPages := currentOffset + pageCount*16
	var pvfPages [][]byte
	var totalPayload uint64

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != fileExt {
			continue
		}

		pagePath := filepath.Join(tempDir, file.Name())
		pageData, err := os.ReadFile(pagePath)
		if err != nil {
			return fmt.Errorf("failed to read split page %s: %w", file.Name(), err)
		}
		if err := os.Remove(pagePath); err != nil {
			return fmt.Errorf("failed to remove page %s: %w", file.Name(), err)
		}

		pvfPages = append(pvfPages, pageData)
		binary.Write(outputFile, binary.LittleEndian, offsetForPages)
		binary.Write(outputFile, binary.LittleEndian, uint64(len(pageData)))
		currentOffset += 16
		offsetForPages += uint64(len(pageData))
		totalPayload += uint64(len(pageData))
	}

	for _, page := range pvfPages {
		outputFile.Write(page)
		currentOffset += uint64(len(page))
	}
	if currentOffset != offsetForPages {
		return fmt.Errorf("wrong offsets for pvf file. Expected %d, got %d", offsetForPages, currentOffset)
	}

	log.Infof("writePVFByFileEntries done ext=%s pages=%d payload=%dB out=%s took=%s", fileExt, pageCount, totalPayload, outputFilePath, time.Since(start))
	return nil
}

func getNum(name string) int {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	parts := strings.Split(base, "_")
	n, _ := strconv.Atoi(parts[len(parts)-1])
	return n
}

func countEntriesByExt(files []os.DirEntry, ext string) int {
	count := 0
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ext {
			continue
		}
		count++
	}
	return count
}
