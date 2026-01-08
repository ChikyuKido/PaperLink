package entity

type Digi4SchoolBook struct {
	ID        int `gorm:"primaryKey"`
	BookName  string
	BookID    string `gorm:"unique"`
	AccountID int
	Account   Digi4SchoolAccount `gorm:"foreignKey:AccountID"`
}
