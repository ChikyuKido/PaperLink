package auth

import (
	"net/http"
	"paperlink/db"
	"paperlink/db/entity"
	"paperlink/server/routes"
	"paperlink/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var log = util.GroupLog("AUTH")

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
		log.WithError(err).Warn("invalid register body")
		c.JSON(http.StatusBadRequest, routes.NewError(400, "invalid request body"))
		return
	}

	// Invite-Code check
	if req.InviteCode != "test" {
		c.JSON(http.StatusForbidden, routes.NewError(403, "invalid invite code"))
		return
	}

	// Existiert Name schon?
	var existing entity.User
	err := db.DB.Where("name = ?", req.Username).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, routes.NewError(409, "username already taken"))
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithError(err).Error("failed to check existing user")
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "internal error"))
		return
	}

	// Passwort hashen
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("failed to hash password")
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "internal error"))
		return
	}

	user := entity.User{
		Name:         req.Username,
		PasswordHash: string(hash),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		log.WithError(err).Error("failed to create user")
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "failed to create user"))
		return
	}

	// JWT erzeugen
	token, err := util.GenerateJWT(user.ID, user.Name)
	if err != nil {
		log.WithError(err).Error("failed to generate jwt")
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "failed to generate jwt"))
		return
	}

	c.JSON(http.StatusOK, routes.NewSuccess(gin.H{
		"jwt": token,
	}))
}
