package entity

type Digi4SchoolAccount struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `gorm:"size:255" json:"password"`
}
