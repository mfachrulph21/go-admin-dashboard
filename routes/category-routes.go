package routes

import (
	"admin-dashboard-FP/handlers"
	"admin-dashboard-FP/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoryRoutes(router *gin.Engine, db *gorm.DB) {
	categoryRouter := router.Group("/categories")
	{
		categoryRouter.Use(middleware.AuthenMiddleware())
		categoryRouter.Use(middleware.AuthorAdminMiddleware())

		categoryRouter.GET("", handlers.GetCategories(db))
		categoryRouter.POST("", handlers.CreateCategory(db))
		categoryRouter.PATCH("/:categoryId", handlers.EditCategory(db))
		categoryRouter.DELETE("/:categoryId", handlers.DeleteCategory(db))
	}
}
