package http

import (
	"github.com/MishaZhem/Gona/src/gateway/internal/app"
	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.POST("/register", register(a))
	r.POST("/login", login(a))
	r.GET("/validate", validateToken(a))
}
