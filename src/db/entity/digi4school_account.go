package entity

type Digi4SchoolAccount struct {
	ID       int    `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string `gorm:"size:255"`
}
