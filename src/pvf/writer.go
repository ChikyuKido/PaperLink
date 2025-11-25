package pvf

import (
	"encoding/binary"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var PVF_MAGIC_BYTES = []byte{0x50, 0x56, 0x46, 0x0A}
var VERSION = []byte{0x1}

func WritePVFFromPDF(inputFile string) (string, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "pvf_*")
	fmt.Println(tempDir)
	if err != nil {
		return "", fmt.Errorf("could not create temporary directory: %w", err)
	}
	//defer os.RemoveAll(tempDir)

	err = api.SplitFile(inputFile, tempDir, 1, nil)
	if err != nil {
		return "", fmt.Errorf("pdfcpu failed to split file: %w", err)
	}
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
	outputFilePath := fmt.Sprintf("%s/output.pvf", tempDir)
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return "", fmt.Errorf("could not create output file: %w", err)
	}

	pageCount := uint64(len(files))
	outputFile.Write(PVF_MAGIC_BYTES)                        // 4
	outputFile.Write(VERSION)                                // 1
	binary.Write(outputFile, binary.LittleEndian, pageCount) // 8
	currentOffset := uint64(13)
	offsetForPages := currentOffset + pageCount*8
	var pvfPages [][]byte
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".pdf" {
			continue
		}

		pagePath := filepath.Join(tempDir, file.Name())

		pageData, err := os.ReadFile(pagePath)
		if err != nil {
			return "", fmt.Errorf("failed to read split page %s: %w", file.Name(), err)
		}

		pvfPages = append(pvfPages, pageData)
		binary.Write(outputFile, binary.LittleEndian, offsetForPages) // 8
		currentOffset += 8
		offsetForPages += uint64(len(pageData))
	}
	for _, page := range pvfPages {
		outputFile.Write(page)
		currentOffset += uint64(len(page))
	}
	if currentOffset != offsetForPages {
		return "", fmt.Errorf("wrong offsets for pvf file. Expected %d, got %d", offsetForPages, currentOffset)
	}

	return "test", nil
}
