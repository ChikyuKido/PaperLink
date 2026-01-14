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

type LoginResponse struct {
	AccessToken string `json:"access"`
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates a user and returns a JWT access token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest        true  "Login payload"
// @Success 200 {object} LoginResponse
// @Failure      400      {object}  routes.ErrorResponse "Invalid request body"
// @Failure      401      {object}  routes.ErrorResponse "Invalid credentials"
// @Failure      500      {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("invalid login body: %v", err)
		routes.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := repo.User.GetUserByName(req.Username)
	if err != nil {
		log.Warnf("login failed for user %s: %v", req.Username, err)
		routes.JSONError(c, http.StatusUnauthorized, "wrong username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Warnf("password mismatch for user %s", req.Username)
		routes.JSONError(c, http.StatusUnauthorized, "wrong username or password")
		return
	}

	access, refresh, err := util.GenerateJWT(user.ID, user.Username)
	if err != nil {
		log.Errorf("failed to generate jwt for user %s: %v", req.Username, err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to generate jwt")
		return
	}

	c.SetCookie(
		"refresh",
		refresh,
		60*60*24*30,
		"/",
		"",
		false,
		true,
	)
	routes.JSONSuccess(c, http.StatusOK, LoginResponse{
		AccessToken: access,
	})

}
