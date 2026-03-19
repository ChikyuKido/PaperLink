package helper

import (
	"encoding/binary"
	"fmt"
	"os"
)

var pvfMagicBytes = []byte{0x50, 0x56, 0x46, 0x0A}
var pvfVersion = []byte{0x1}

func ConvertPDFPagesToPVF(pagePDFs []string, outputPVF string) error {
	outputFile, err := os.Create(outputPVF)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer outputFile.Close()

	pageCount := uint64(len(pagePDFs))
	if pageCount == 0 {
		return fmt.Errorf("no page pdfs provided")
	}

	if _, err := outputFile.Write(pvfMagicBytes); err != nil {
		return fmt.Errorf("failed to write pvf magic: %w", err)
	}
	if _, err := outputFile.Write(pvfVersion); err != nil {
		return fmt.Errorf("failed to write pvf version: %w", err)
	}
	if err := binary.Write(outputFile, binary.LittleEndian, pageCount); err != nil {
		return fmt.Errorf("failed to write page count: %w", err)
	}

	currentOffset := uint64(13)
	offsetForPages := currentOffset + pageCount*16
	pages := make([][]byte, 0, pageCount)

	for _, pagePath := range pagePDFs {
		pageData, err := os.ReadFile(pagePath)
		if err != nil {
			return fmt.Errorf("failed to read split page %s: %w", pagePath, err)
		}
		pages = append(pages, pageData)

		if err := binary.Write(outputFile, binary.LittleEndian, offsetForPages); err != nil {
			return fmt.Errorf("failed to write page offset: %w", err)
		}
		if err := binary.Write(outputFile, binary.LittleEndian, uint64(len(pageData))); err != nil {
			return fmt.Errorf("failed to write page size: %w", err)
		}
		offsetForPages += uint64(len(pageData))
	}

	for _, pageData := range pages {
		if _, err := outputFile.Write(pageData); err != nil {
			return fmt.Errorf("failed to write page data: %w", err)
		}
	}

	return nil
}
