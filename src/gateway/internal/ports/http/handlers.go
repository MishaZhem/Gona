package http

import (
	"errors"
	"net/http"

	"github.com/MishaZhem/Gona/src/gateway/internal/app"
	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

var ErrInvalid = errors.New("invalid request")
var ErrNoToken = errors.New("missing token")

func register(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(ErrInvalid))
			return
		}

		err := a.Register(req.Username, req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		c.Status(http.StatusCreated)
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func login(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(ErrInvalid))
			return
		}

		token, err := a.Login(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse(err))
			return
		}

		c.SetCookie(
			"access_token",
			token,
			3600,
			"/",
			"",
			true,
			true,
		)

		c.JSON(http.StatusOK, gin.H{"message": "login successful"})
	}
}

func logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie(
			"access_token",
			"",
			-1,
			"/",
			"",
			true,
			true,
		)

		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	}
}

func userProfile(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID missing"})
			return
		}

		profile, err := a.UserProfile(userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
			return
		}

		c.JSON(http.StatusOK, profile)
	}
}
