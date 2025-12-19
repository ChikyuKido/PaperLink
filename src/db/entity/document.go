package entity

type Document struct {
	ID          int    `gorm:"primaryKey"`
	UUID        string `gorm:"uniqueIndex;not null"`
	Name        string
	Description string

	Tags []Tag `gorm:"many2many:document_tags;constraint:OnDelete:CASCADE"`

	UserID int
	User   User `gorm:"constraint:OnDelete:CASCADE"`

	FileUUID string
	File     FileDocument `gorm:"foreignKey:FileUUID;references:UUID;constraint:OnDelete:CASCADE"`

	DirectoryID *int
	Directory   *Directory `gorm:"foreignKey:DirectoryID;constraint:OnDelete:CASCADE"`
}
