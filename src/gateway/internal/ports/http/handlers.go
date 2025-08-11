package http

import (
	"errors"
	"io"
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

func uploadAvatar(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID missing"})
			return
		}

		file, _, err := c.Request.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}
		fileSize := int64(len(fileBytes))

		probe := fileBytes
		if len(probe) > 512 {
			probe = probe[:512]
		}
		contentType := http.DetectContentType(probe)

		allowed := map[string]bool{
			"image/png":  true,
			"image/jpeg": true,
			"image/webp": true,
		}
		if !allowed[contentType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
			return
		}

		url, err := a.ChangeAvatar(userID.(string), fileBytes, fileSize, contentType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}

		c.JSON(http.StatusOK, url)
	}
}

type updateStatusWorkerRequest struct {
	Worker bool `json:"worker" binding:"required"`
}

func updateStatusWorker(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updateStatusWorkerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(ErrInvalid))
			return
		}

		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID missing"})
			return
		}

		err := a.UpdateStatusWorker(userID.(string), req.Worker)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		c.Status(http.StatusCreated)
	}
}
