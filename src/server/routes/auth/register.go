package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"
	"time"
)

type RegisterRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	InviteCode string `json:"inviteCode"`
}

// Register godoc
// @Summary      Register user
// @Description  Creates a new user using a valid invite code.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest true "Register payload"
// @Success      200      {object}  routes.Response
// @Failure      400      {object}  routes.ErrorResponse "Invalid request body"
// @Failure      401      {object}  routes.ErrorResponse "Invalid invite code"
// @Failure      409      {object}  routes.ErrorResponse "Username already taken"
// @Failure      500      {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/auth/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("invalid register body: %v", err)
		routes.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	invite, err := repo.RegistrationInvite.GetByCode(req.InviteCode)
	if err != nil || invite == nil {
		routes.JSONError(c, http.StatusUnauthorized, "invite code invalid")
		return
	}
	if invite.ExpiresAt < time.Now().Unix() || invite.Uses <= 0 {
		err := repo.RegistrationInvite.Delete(invite.ID)
		if err != nil {
			log.Warnf("failed to delete expired invite: %v", err)
		}
		routes.JSONError(c, http.StatusUnauthorized, "invite expired")
		return
	}

	existingUser, err := repo.User.GetUserByName(req.Username)
	if err == nil && existingUser != nil && existingUser.ID != 0 {
		routes.JSONError(c, http.StatusConflict, "username already taken")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("failed to hash password: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to create user")
		return
	}

	user := entity.User{
		Username: req.Username,
		Password: string(hash),
		IsAdmin:  req.InviteCode == "admin",
	}

	if err := repo.User.Save(&user); err != nil {
		log.Errorf("failed to create user: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to create user")
		return
	}

	invite.Uses -= 1
	if err := repo.RegistrationInvite.Save(invite); err != nil {
		log.Warnf("failed to update invite: %v", err)
	}

	routes.JSONSuccess(c, http.StatusOK, gin.H{
		"message": "ok",
	})
}
