package repo

import (
	"crypto/rand"
	"encoding/hex"
	"paperlink/db/entity"
	"time"
)

type RegistrationInviteRepo struct {
	*Repository[entity.RegistrationInvite]
}

func newRegistrationInviteRepo() *RegistrationInviteRepo {
	return &RegistrationInviteRepo{NewRepository[entity.RegistrationInvite]()}
}

// Global nutzbares Repo-Objekt
var RegistrationInvite = newRegistrationInviteRepo()

// intern: Random Code generieren
func generateInviteCode() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// Create erzeugt einen neuen Invite-Code (default: 3 Tage gültig)
func (r *RegistrationInviteRepo) Create(validDays int) (*entity.RegistrationInvite, error) {
	if validDays <= 0 {
		validDays = 3
	}

	code, err := generateInviteCode()
	if err != nil {
		return nil, err
	}

	invite := entity.RegistrationInvite{
		Code:      code,
		ExpiresAt: time.Now().Add(time.Duration(validDays) * 24 * time.Hour).Unix(),
	}

	if err := r.Save(&invite); err != nil {
		return nil, err
	}

	return &invite, nil
}

// GetByCode holt einen Invite anhand des Codes
func (r *RegistrationInviteRepo) GetByCode(code string) (*entity.RegistrationInvite, error) {
	var invite entity.RegistrationInvite
	if err := r.db.Where("code = ?", code).First(&invite).Error; err != nil {
		return nil, err
	}
	return &invite, nil
}
