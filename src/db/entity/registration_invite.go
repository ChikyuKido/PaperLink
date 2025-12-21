package entity

type RegistrationInvite struct {
	ID        int    `gorm:"primary_key;AUTO_INCREMENT"`
	Code      string `gorm:"uniqueIndex;not null"`
	Uses      int    `gorm:"not null"`
	ExpiresAt int64  `gorm:"not null"`
}
