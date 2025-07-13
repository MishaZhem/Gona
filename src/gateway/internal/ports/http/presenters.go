package http

import (
	"github.com/gin-gonic/gin"
)

func TokenResponse(token string) gin.H {
	return gin.H{
		"token": token,
	}
}

func UserIdResponse(userID string) gin.H {
	return gin.H{
		"user_id": userID,
	}
}

func ErrorResponse(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
