package entity

type Document struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT"`
	UUID         string `gorm:"uniqueIndex"`
	FileUUID     string
	FileDocument FileDocument `gorm:"foreignKey:FileUUID"`
	Name         string
	Description  string
	DirectoryID  *int
	Directory    *Directory `gorm:"foreignKey:DirectoryID;references:ID;constraint:OnDelete:CASCADE"`
	User         User
	UserID       int
	Tags         []Tag `gorm:"many2many:document_tags"`
}
