package routes

import (
	"admin-dashboard-FP/handlers"
	"admin-dashboard-FP/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProductRoutes(router *gin.Engine, db *gorm.DB) {
	productRouter := router.Group("/products")
	{
		productRouter.Use(middleware.AuthenMiddleware())

		productRouter.GET("", handlers.GetProducts(db))
		productRouter.POST("", middleware.AuthorAdminMiddleware(), handlers.CreateProduct(db))
		productRouter.PUT("/:productId", middleware.AuthorAdminMiddleware(), handlers.EditProduct(db))
		productRouter.DELETE("/:productId", middleware.AuthorAdminMiddleware(), handlers.DeleteProduct(db))
	}
}
