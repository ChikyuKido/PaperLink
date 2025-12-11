package auth

import (
	"net/http"
	"paperlink/db/repo"
	"paperlink/server/routes"
	"paperlink/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	user, err := repo.User.GetUserByName(req.Username)
	if err != nil {
		log.Errorf("failed to query user: %v", err)
		c.JSON(http.StatusUnauthorized, routes.NewError(401, "wrong username or password"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, routes.NewError(401, "wrong username or password"))
		return
	}

	token, err := util.GenerateJWT(user.ID, user.Username)
	if err != nil {
		log.Errorf("failed to generate jwt: %v", err)
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "failed to generate jwt"))
		return
	}

	c.JSON(http.StatusOK, routes.NewSuccess(gin.H{
		"jwt": token,
	}))
}
