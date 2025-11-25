package entity

type DocumentUserRole string

const (
	Editor DocumentUserRole = "EDITOR"
	Viewer DocumentUserRole = "VIEWER"
)

type DocumentUser struct {
	User       User
	UserId     int `gorm:"foreignKey:UserID;primaryKey"`
	Document   Document
	DocumentID int `gorm:"foreignKey:DocumentID;primaryKey"`
	Role       DocumentUserRole
}
