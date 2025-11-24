package entity

type Document struct {
	ID         int `gorm:"primary_key;AUTO_INCREMENT"`
	Name       string
	Size       uint64
	Pages      uint64
	Category   Category
	CategoryID int   `gorm:"foreignKey:CategoryID"`
	Tags       []Tag `gorm:"many2many:document_tags"`
}
