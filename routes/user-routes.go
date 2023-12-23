package routes

import (
	"admin-dashboard-FP/handlers"
	"admin-dashboard-FP/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(router *gin.Engine, db *gorm.DB) {
	userRouter := router.Group("/users")
	{
		userRouter.POST("/register", handlers.UserRegister(db))
		userRouter.POST("/login", handlers.UserLogin(db))
		userRouter.PATCH("/topup", middleware.AuthenMiddleware(), handlers.UserTopup(db))
	}
}
