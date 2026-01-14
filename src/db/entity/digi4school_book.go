package entity

type Digi4SchoolBook struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	UUID      string `gorm:"unique" json:"uuid"`
	BookName  string `json:"bookName"`
	BookID    string `gorm:"unique" json:"bookId"`
	AccountID int
	Account   Digi4SchoolAccount `gorm:"foreignKey:AccountID" json:"account"`

	FileUUID string
	File     FileDocument `gorm:"foreignKey:FileUUID;references:UUID" json:"file"`
}
