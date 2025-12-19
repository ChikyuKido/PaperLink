package entity

type Directory struct {
	ID       int        `gorm:"primaryKey"`
	UserID   int        `gorm:"not null;index:idx_user_parent_name,unique"`
	User     User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ParentID *int       `gorm:"index:idx_user_parent_name,unique"`
	Parent   *Directory `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
	Name     string     `gorm:"size:255;not null;index:idx_user_parent_name,unique"`
}
