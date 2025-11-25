package entity

type Document struct {
	ID       int `gorm:"primary_key;AUTO_INCREMENT"`
	UUID     string
	FileUUID string `gorm:"foreignKey:FileDocumentUUID"`
	Name     string
	Size     uint64
	Pages    uint64
	Owner    User
	OwnerID  int
	Tags     []Tag `gorm:"many2many:document_tags"`
}
