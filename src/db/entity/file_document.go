package entity

type FileDocument struct {
	UUID string `gorm:"primary_key"`
	Path string `gorm:"uniqueIndex"`
}
