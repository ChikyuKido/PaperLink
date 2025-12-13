package entity

type Document struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT"`
	UUID        string `gorm:"uniqueIndex"`
	FileUUID    string `gorm:"foreignKey:FileDocumentUUID"`
	Name        string
	Description string
	Path        string
	Owner       User
	OwnerID     int
	Tags        []Tag `gorm:"many2many:document_tags"`
}
