package routes

import (
	"admin-dashboard-FP/handlers"
	"admin-dashboard-FP/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransactionRoutes(router *gin.Engine, db *gorm.DB) {
	transactionRoutes := router.Group("/transactions")
	{
		transactionRoutes.Use(middleware.AuthenMiddleware())

		transactionRoutes.GET("/my-transactions", handlers.GetTransactionLoginUser(db))
		transactionRoutes.GET("/user-transactions", middleware.AuthorAdminMiddleware(), handlers.GetTransactionUserByAdmin(db))
		transactionRoutes.POST("", handlers.PostTransaction(db))
	}
}
