package entity

type Document struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT"`
	UUID        string `gorm:"uniqueIndex"`
	FileUUID    string `gorm:"foreignKey:FileDocumentUUID"`
	Name        string
	Description string
	DirectoryID *int
	Directory   *Directory `gorm:"foreignKey:DirectoryID;references:ID;constraint:OnDelete:CASCADE"`
	Owner       User
	OwnerID     int
	Tags        []Tag `gorm:"many2many:document_tags"`
}
