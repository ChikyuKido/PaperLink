package entity

type User struct {
	ID       int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name     string `gorm:"unique;not null"`
	Password string
	IsAdmin  bool
}
