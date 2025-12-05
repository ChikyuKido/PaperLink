package auth

import (
	"net/http"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	InviteCode string `json:"inviteCode"`
}

// Register godoc
// @Summary      Register new user
// @Description  Creates a new user. Invite code must be "test".
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Register payload"
// @Success      200      {object}  routes.Response
// @Failure      400      {object}  routes.Response
// @Failure      403      {object}  routes.Response
// @Failure      409      {object}  routes.Response
// @Failure      500      {object}  routes.Response
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("invalid register body: %v", err)
		c.JSON(http.StatusBadRequest, routes.NewError(400, "invalid request body"))
		return
	}

	if req.InviteCode != "test" {
		c.JSON(http.StatusForbidden, routes.NewError(403, "invalid invite code"))
		return
	}

	exists := repo.User.DoesUserByNameExist(req.Username)
	if exists {
		c.JSON(http.StatusConflict, routes.NewError(409, "username already taken"))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("failed to hash password")
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "internal error"))
		return
	}

	user := entity.User{
		Name:     req.Username,
		Password: string(hash),
		IsAdmin:  false,
	}

	if err := repo.User.Save(&user).Error; err != nil {
		log.Errorf("failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "failed to create user"))
		return
	}

	c.JSON(http.StatusOK, routes.NewSuccess(gin.H{
		"message": "ok",
	}))
}
