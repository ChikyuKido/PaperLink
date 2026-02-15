package ptf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

type ReadOptions struct {
	HasRange bool
	Start    uint64
	End      uint64
}

func Read(filePath string, opts ReadOptions) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if !opts.HasRange {
		return io.ReadAll(file)
	}
	if opts.End < opts.Start {
		return nil, fmt.Errorf("invalid range: end must be >= start")
	}

	entries, err := readIndexTable(file)
	if err != nil {
		return nil, err
	}
	return readRangeByIndex(file, entries, opts.Start, opts.End)
}

func readIndexTable(file *os.File) ([]pageEntry, error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	magic := make([]byte, headerMagicSize)
	if _, err := io.ReadFull(file, magic); err != nil {
		return nil, fmt.Errorf("failed to read ptf magic: %w", err)
	}

	if !bytes.Equal(magic, fileMagic[:]) {
		return nil, fmt.Errorf("invalid ptf magic")
	}

	header := make([]byte, headerVersionSize+headerMapSizeFieldSize)
	if _, err := io.ReadFull(file, header); err != nil {
		return nil, fmt.Errorf("failed to read ptf header: %w", err)
	}

	version := header[0]
	if version != fileVersionIndexed {
		return nil, fmt.Errorf("unsupported ptf version: %d", version)
	}

	mapSize := binary.LittleEndian.Uint64(header[1:])
	if mapSize < 8 {
		return nil, fmt.Errorf("invalid ptf map size: %d", mapSize)
	}
	if mapSize > math.MaxInt64 {
		return nil, fmt.Errorf("ptf map too large: %d", mapSize)
	}

	mapBytes := make([]byte, mapSize)
	if _, err := io.ReadFull(file, mapBytes); err != nil {
		return nil, fmt.Errorf("failed to read ptf index map: %w", err)
	}

	pageCount := binary.LittleEndian.Uint64(mapBytes[:8])
	expectedMapSize := uint64(8) + pageCount*indexEntrySize
	if expectedMapSize != mapSize {
		return nil, fmt.Errorf("invalid ptf map layout: mapSize=%d expected=%d", mapSize, expectedMapSize)
	}

	entries := make([]pageEntry, pageCount)
	offset := 8
	for i := uint64(0); i < pageCount; i++ {
		entries[i] = pageEntry{
			Offset: binary.LittleEndian.Uint64(mapBytes[offset : offset+8]),
			Size:   binary.LittleEndian.Uint64(mapBytes[offset+8 : offset+16]),
		}
		offset += indexEntrySize
	}

	return entries, nil
}

func readRangeByIndex(file *os.File, entries []pageEntry, start, end uint64) ([]byte, error) {
	pageCount := uint64(len(entries))
	if start >= pageCount {
		return nil, os.ErrNotExist
	}
	if end >= pageCount {
		end = pageCount - 1
	}

	var out bytes.Buffer
	for i := start; i <= end; i++ {
		entry := entries[i]
		if entry.Size > math.MaxInt64 {
			return nil, fmt.Errorf("thumbnail page size too large: %d", entry.Size)
		}
		if entry.Offset > math.MaxInt64 {
			return nil, fmt.Errorf("thumbnail page offset too large: %d", entry.Offset)
		}

		buf := make([]byte, entry.Size)
		n, err := file.ReadAt(buf, int64(entry.Offset))
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n != int(entry.Size) {
			return nil, io.ErrUnexpectedEOF
		}
		if err := binary.Write(&out, binary.LittleEndian, entry.Size); err != nil {
			return nil, err
		}
		if _, err := out.Write(buf); err != nil {
			return nil, err
		}
	}

	return out.Bytes(), nil
}
