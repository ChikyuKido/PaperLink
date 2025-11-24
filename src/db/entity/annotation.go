package entity

type AnnotationType string

const (
	Textbox AnnotationType = "TEXTBOX"
	Note    AnnotationType = "NOTE"
	Canvas  AnnotationType = "CANVAS"
)

type Annotation struct {
	ID         int `gorm:"primary_key;AUTO_INCREMENT"`
	Type       AnnotationType
	CreatedAt  int64
	UpdatedAt  int64
	PositionX  float64
	PositionY  float64
	Document   Document
	DocumentID int `gorm:"foreignKey:DocumentID"`
}
