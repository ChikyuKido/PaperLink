package entity

type DocumentUserRole string

const (
	Editor DocumentUserRole = "EDITOR"
	Viewer DocumentUserRole = "VIEWER"
)

type DocumentUser struct {
	UserID     int      `gorm:"primaryKey"`
	DocumentID int      `gorm:"primaryKey"`
	User       User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Document   Document `gorm:"foreignKey:DocumentID;constraint:OnDelete:CASCADE"`
	Role       DocumentUserRole
}
