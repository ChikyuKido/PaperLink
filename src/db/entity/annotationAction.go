package entity

type Action string

const (
	Create Action = "CREATE"
	Update Action = "UPDATE"
	Move   Action = "MOVE"
	Delete Action = "DELETE"
)

type AnnotationAction struct {
	ID           int `gorm:"primary_key;AUTO_INCREMENT"`
	Data         string
	CreatedAt    int64
	Annotation   Annotation
	AnnotationID int `gorm:"foreignKey:AnnotationID"`
}
