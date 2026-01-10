package helper

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

var PVF_MAGIC_BYTES = []byte{0x50, 0x56, 0x46, 0x0A}
var VERSION = []byte{0x1}

func WritePVFFromPDF(files []string, outputFilePath string) error {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}

	pageCount := uint64(len(files))
	outputFile.Write(PVF_MAGIC_BYTES)                        // 4
	outputFile.Write(VERSION)                                // 1
	binary.Write(outputFile, binary.LittleEndian, pageCount) // 8
	currentOffset := uint64(13)
	offsetForPages := currentOffset + pageCount*16
	var pvfPages [][]byte
	for _, filePath := range files {
		file, err := os.Stat(filePath)
		if file.IsDir() || filepath.Ext(file.Name()) != ".pdf" {
			continue
		}

		pageData, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read split page %s: %w", file.Name(), err)
		}
		pvfPages = append(pvfPages, pageData)
		binary.Write(outputFile, binary.LittleEndian, offsetForPages)        // 8
		binary.Write(outputFile, binary.LittleEndian, uint64(len(pageData))) // 8
		currentOffset += 16
		offsetForPages += uint64(len(pageData))
	}
	for _, page := range pvfPages {
		outputFile.Write(page)
		currentOffset += uint64(len(page))
	}
	if currentOffset != offsetForPages {
		return fmt.Errorf("wrong offsets for pvf file. Expected %d, got %d", offsetForPages, currentOffset)
	}

	return nil
}
