package entity

type User struct {
	ID      int `gorm:"primary_key;AUTO_INCREMENT"`
	Name    string
	IsAdmin bool
}
