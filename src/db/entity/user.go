package entity

type User struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"uniqueIndex"`
	PasswordHash string
	IsAdmin      bool
}
