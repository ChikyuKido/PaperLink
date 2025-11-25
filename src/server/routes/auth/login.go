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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates a user and returns a JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login payload"
// @Success      200      {object}  routes.Response
// @Failure      400      {object}  routes.Response
// @Failure      401      {object}  routes.Response
// @Failure      500      {object}  routes.Response
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("invalid login body")
		c.JSON(http.StatusBadRequest, routes.NewError(400, "invalid request body"))
		return
	}

	var user entity.User
	if err := db.DB.Where("name = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, routes.NewError(401, "wrong username or password"))
			return
		}
		log.WithError(err).Error("failed to query user")
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "internal error"))
		return
	}

	// Passwort prüfen
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, routes.NewError(401, "wrong username or password"))
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
