package auth

import (
	"net/http"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	InviteCode string `json:"inviteCode"`
}

// POST /api/v1/auth/register
// erwartet gültigen Invite-Code
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("invalid register body: %v", err)
		c.JSON(http.StatusBadRequest, routes.NewError(400, "invalid request body"))
		return
	}

	// Invite-Code prüfen
	invite, err := repo.RegistrationInvite.GetByCode(req.InviteCode)
	if err != nil || invite == nil || invite.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, routes.NewError(401, "invite code invalid"))
		return
	}

	// prüfen, ob Username schon existiert
	existingUser, err := repo.User.GetUserByName(req.Username)
	if err == nil && existingUser != nil && existingUser.ID != 0 {
		c.JSON(http.StatusConflict, routes.NewError(409, "username already taken"))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "failed to create user"))
		return
	}

	user := entity.User{
		Username: req.Username,
		Password: string(hash),
		IsAdmin:  false,
	}

	if err := repo.User.Save(&user); err != nil {
		log.Errorf("failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "failed to create user"))
		return
	}

	c.JSON(http.StatusOK, routes.NewSuccess(gin.H{
		"message": "ok",
	}))
}
