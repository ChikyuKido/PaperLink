package entity

type Directory struct {
	ID       int    `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	UserID   int
	ParentID *int

	User   User       `gorm:"constraint:OnDelete:CASCADE"`
	Parent *Directory `gorm:"foreignKey:ParentID"`
}
