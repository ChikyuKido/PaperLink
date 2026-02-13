package entity

type FileDocument struct {
	UUID  string `gorm:"primary_key" json:"uuid"`
	Path  string `gorm:"uniqueIndex" json:"path"`
	Size  uint64 `json:"size"`
	Pages uint64 `json:"pages"`
}
