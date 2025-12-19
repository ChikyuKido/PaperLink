package entity

type Directory struct {
	ID       int    `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	UserID   int
	ParentID *int

	// delete all directories if user exists
	User User `gorm:"constraint:OnDelete:CASCADE"`

	Parent *Directory `gorm:"foreignKey:ParentID"`
}
