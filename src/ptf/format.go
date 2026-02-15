package ptf

const (
	headerMagicSize        = 4
	headerVersionSize      = 1
	headerMapSizeFieldSize = 8
	headerFixedSize        = headerMagicSize + headerVersionSize + headerMapSizeFieldSize
	indexEntrySize         = 16
)

var fileMagic = [headerMagicSize]byte{0x50, 0x54, 0x46, 0x0A}

const fileVersionIndexed byte = 0x1

type pageEntry struct {
	Offset uint64
	Size   uint64
}
