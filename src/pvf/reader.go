package pvf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type Metadata struct {
	PageCount uint64
	Indexes   []PageData
}
type PageData struct {
	Offset uint64
	Size   uint64
}

func ReadMetadata(filePath string) (Metadata, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return Metadata{}, fmt.Errorf("filePath %s does not exist", filePath)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to open file %s", filePath)
	}
	data := make([]byte, 4)
	read, err := file.Read(data)
	if err != nil || read != 4 || !bytes.Equal(data, PVF_MAGIC_BYTES) {
		return Metadata{}, fmt.Errorf("failed to check magic bytes of the file: %s", filePath)
	}
	data = make([]byte, 1)
	read, err = file.Read(data)
	if err != nil || read != 1 {
		return Metadata{}, fmt.Errorf("failed to read version from file: %s", filePath)
	}
	var pageCount uint64
	err = binary.Read(file, binary.LittleEndian, &pageCount)
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to read page count from file: %s", filePath)
	}
	indexes := make([]PageData, pageCount)
	for i := uint64(0); i < pageCount; i++ {
		var indexOffset uint64
		var size uint64
		err = binary.Read(file, binary.LittleEndian, &indexOffset)
		if err != nil {
			return Metadata{}, fmt.Errorf("failed to read page index from file: %s", filePath)
		}
		err = binary.Read(file, binary.LittleEndian, &size)
		if err != nil {
			return Metadata{}, fmt.Errorf("failed to read page index from file: %s", filePath)
		}
		indexes[i].Size = size
		indexes[i].Offset = indexOffset
	}

	return Metadata{
		PageCount: pageCount,
		Indexes:   indexes,
	}, nil
}

func ReadPage(filePath string, page uint64) ([]byte, error) {
	metadata, err := ReadMetadata(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read page metadata: %s", err)
	}
	if page == 0 {
		return nil, fmt.Errorf("page number is zero")
	}
	page -= 1
	pageData := metadata.Indexes[page]
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s", filePath)
	}
	data := make([]byte, pageData.Size)
	read, err := file.ReadAt(data, int64(pageData.Offset))
	if err != nil || read != int(pageData.Size) {
		return nil, fmt.Errorf("failed to read page from file: %s", filePath)
	}
	return data, nil
}

func ReadPages(filePath string, start, end uint64) ([][]byte, error) {
	if start == 0 || end == 0 {
		return nil, fmt.Errorf("page numbers must be >= 1")
	}
	if end < start {
		return nil, fmt.Errorf("end page must be >= start page")
	}

	md, err := ReadMetadata(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %s", err)
	}

	if start > md.PageCount || end > md.PageCount {
		return nil, fmt.Errorf("page range out of bounds")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %s", filePath, err)
	}
	defer file.Close()

	result := make([][]byte, 0, end-start+1)

	for p := start; p <= end; p++ {
		idx := p - 1
		pageData := md.Indexes[idx]

		buf := make([]byte, pageData.Size)
		n, err := file.ReadAt(buf, int64(pageData.Offset))
		if err != nil || n != int(pageData.Size) {
			return nil, fmt.Errorf("failed to read page %d: %s", p, err)
		}
		result = append(result, buf)
	}

	return result, nil
}
