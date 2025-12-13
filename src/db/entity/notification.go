package entity

type ActionType string

const (
	Invite ActionType = "INVITE"
)

type Notification struct {
	ID         int `gorm:"primary_key;AUTO_INCREMENT"`
	Title      string
	Text       string
	Action     ActionType
	ActionData string
	User       User
	UserID     int `gorm:"foreignKey:UserID"`
}
