package http

import (
	"github.com/MishaZhem/Gona/src/gateway/internal/app"
	"github.com/MishaZhem/Gona/src/gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.POST("/register", register(a))
	r.POST("/login", login(a))
	r.POST("/logout", logout())

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(a))
	auth.GET("/profile", userProfile(a))
	auth.POST("/upload-avatar", uploadAvatar(a))
	auth.POST("/update-status-worker", updateStatusWorker(a))

	auth.GET("/ws/support")
	auth.GET("/ws/support/agent")
}
