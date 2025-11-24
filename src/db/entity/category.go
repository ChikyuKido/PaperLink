package entity

type Category struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	Name     string
	ParentID *int
	Parent   *Category `gorm:"foreignKey:ParentID"`
}
