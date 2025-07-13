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

		c.JSON(http.StatusOK, TokenResponse(token))
	}
}

func validateToken(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse(ErrNoToken))
			return
		}

		userID, err := a.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserIdResponse(userID))
	}
}
