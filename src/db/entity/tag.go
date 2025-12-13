package entity

type Tag struct {
	ID    int `gorm:"primary_key;AUTO_INCREMENT"`
	Name  string
	Color string
}
