package pvf

import (
	"encoding/binary"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"os"
	"os/exec"
	"paperlink/util"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var PVF_MAGIC_BYTES = []byte{0x50, 0x56, 0x46, 0x0A}
var VERSION = []byte{0x1}
var log = util.GroupLog("PVF")

func WritePVFFromPDF(inputFile string) (string, error) {
	start := time.Now()
	log.Infof("WritePVFFromPDF start input=%s", inputFile)

	tempDir, err := os.MkdirTemp(os.TempDir(), "pvf_*")
	if err != nil {
		return "", fmt.Errorf("could not create temporary directory: %w", err)
	}
	splitStart := time.Now()
	err = api.SplitFile(inputFile, tempDir, 1, nil)
	if err != nil {
		return "", fmt.Errorf("pdfcpu failed to split file: %w", err)
	}
	log.Infof("WritePVFFromPDF split done dir=%s took=%s", tempDir, time.Since(splitStart))
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return "", fmt.Errorf("could not read temporary directory: %w", err)
	}
	sort.Slice(files, func(i, j int) bool {
		getNum := func(name string) int {
			base := strings.TrimSuffix(name, ".pdf")
			parts := strings.Split(base, "_")
			n, _ := strconv.Atoi(parts[len(parts)-1])
			return n
		}
		return getNum(files[i].Name()) < getNum(files[j].Name())
	})
	pageCount := countEntriesByExt(files, ".pdf")
	outputFilePath := fmt.Sprintf("%s/output.pvf", tempDir)
	writeStart := time.Now()
	err = writePVFByFileEntries(tempDir, files, ".pdf", outputFilePath)
	if err != nil {
		return "", err
	}
	log.Infof("WritePVFFromPDF done pages=%d out=%s write=%s total=%s", pageCount, outputFilePath, time.Since(writeStart), time.Since(start))
	return outputFilePath, nil
}

func WriteThumbnailPTFFromPDF(inputFile string) (string, error) {
	start := time.Now()
	log.Infof("WriteThumbnailPTFFromPDF start input=%s", inputFile)

	tempDir, err := os.MkdirTemp(os.TempDir(), "pvf_thumb_*")
	if err != nil {
		return "", fmt.Errorf("could not create temporary directory: %w", err)
	}

	thumbnailPattern := filepath.Join(tempDir, "thumb_%05d.png")
	cmd := exec.Command(
		"gs",
		"-sDEVICE=pnggray",
		"-r100",
		"-dDownScaleFactor=4",
		"-dTextAlphaBits=4",
		"-dGraphicsAlphaBits=4",
		"-dBATCH",
		"-dNOPAUSE",
		"-sOutputFile="+thumbnailPattern,
		inputFile,
	)
	gsStart := time.Now()
	if out, cmdErr := cmd.CombinedOutput(); cmdErr != nil {
		return "", fmt.Errorf("ghostscript failed to create thumbnails: %w: %s", cmdErr, strings.TrimSpace(string(out)))
	}
	log.Infof("WriteThumbnailPTFFromPDF ghostscript done dir=%s took=%s", tempDir, time.Since(gsStart))

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return "", fmt.Errorf("could not read temporary directory: %w", err)
	}
	sort.Slice(files, func(i, j int) bool {
		getNum := func(name string) int {
			base := strings.TrimSuffix(name, ".png")
			parts := strings.Split(base, "_")
			n, _ := strconv.Atoi(parts[len(parts)-1])
			return n
		}
		return getNum(files[i].Name()) < getNum(files[j].Name())
	})
	pageCount := countEntriesByExt(files, ".png")

	outputFilePath := fmt.Sprintf("%s/output_thumb.ptf", tempDir)
	writeStart := time.Now()
	err = writePTFByFileEntries(tempDir, files, ".png", outputFilePath)
	if err != nil {
		return "", err
	}
	log.Infof("WriteThumbnailPTFFromPDF done pages=%d out=%s write=%s total=%s", pageCount, outputFilePath, time.Since(writeStart), time.Since(start))
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
	outputFile.Write(PVF_MAGIC_BYTES)                        // 4
	outputFile.Write(VERSION)                                // 1
	binary.Write(outputFile, binary.LittleEndian, pageCount) // 8
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
		err = os.Remove(pagePath)
		if err != nil {
			return fmt.Errorf("failed to remove page %s: %w", file.Name(), err)
		}
		pvfPages = append(pvfPages, pageData)
		binary.Write(outputFile, binary.LittleEndian, offsetForPages)        // 8
		binary.Write(outputFile, binary.LittleEndian, uint64(len(pageData))) // 8
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

func writePTFByFileEntries(tempDir string, files []os.DirEntry, fileExt string, outputFilePath string) error {
	start := time.Now()
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer outputFile.Close()

	pageCount := 0
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
		err = os.Remove(pagePath)
		if err != nil {
			return fmt.Errorf("failed to remove page %s: %w", file.Name(), err)
		}

		if err := binary.Write(outputFile, binary.LittleEndian, uint64(len(pageData))); err != nil {
			return fmt.Errorf("failed to write page size for %s: %w", file.Name(), err)
		}
		if _, err := outputFile.Write(pageData); err != nil {
			return fmt.Errorf("failed to write page data for %s: %w", file.Name(), err)
		}
		pageCount++
		totalPayload += uint64(len(pageData))
	}

	log.Infof("writePTFByFileEntries done ext=%s pages=%d payload=%dB out=%s took=%s", fileExt, pageCount, totalPayload, outputFilePath, time.Since(start))
	return nil
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
