package entity

type Digi4SchoolBook struct {
	ID        int    `gorm:"primaryKey"`
	UUID      string `gorm:"unique"`
	BookName  string
	BookID    string `gorm:"unique"`
	AccountID int
	Account   Digi4SchoolAccount `gorm:"foreignKey:AccountID"`

	FileUUID string
	File     FileDocument `gorm:"foreignKey:FileUUID;references:UUID"`
}
