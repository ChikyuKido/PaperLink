package entity

type Document struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT"`
	UUID         string `gorm:"uniqueIndex"`
	FileUUID     string
	FileDocument FileDocument `gorm:"foreignKey:FileUUID;constraint:OnDelete:CASCADE"`
	Name         string
	Description  string
	DirectoryID  *int
	Directory    *Directory `gorm:"foreignKey:DirectoryID;references:ID;constraint:OnDelete:CASCADE"`
	User         User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	UserID       int
	Tags         []Tag `gorm:"many2many:document_tags;constraint:OnDelete:CASCADE"`
}
