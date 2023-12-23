package handlers

import (
	"admin-dashboard-FP/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type resCreateCategory struct {
	ID                uint      `json:"id"`
	Type              string    `json:"type"`
	SoldProductAmount int       `json:"sold_product_amount"`
	CreatedAt         time.Time `json:"created_at"`
}

type resProduct struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type resGetAllCategories struct {
	ID                uint      `json:"id"`
	Type              string    `json:"type"`
	SoldProductAmount int       `json:"sold_product_amount"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Products          []resProduct
}

func GetCategories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category []models.Category

		if err := db.Preload("Products").Find(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Categories not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var categoryResponses []resGetAllCategories

		for _, eachCategory := range category {
			var productResponses []resProduct
			for _, product := range eachCategory.Products {
				productResponse := resProduct{
					ID:        product.ID,
					Title:     product.Title,
					Price:     product.Price,
					Stock:     product.Stock,
					CreatedAt: product.CreatedAt,
					UpdatedAt: product.UpdatedAt,
				}
				productResponses = append(productResponses, productResponse)
			}

			categoryResponse := resGetAllCategories{
				ID:                eachCategory.ID,
				Type:              eachCategory.Type,
				SoldProductAmount: eachCategory.SoldProductAmount,
				CreatedAt:         eachCategory.CreatedAt,
				UpdatedAt:         eachCategory.UpdatedAt,
				Products:          productResponses,
			}

			categoryResponses = append(categoryResponses, categoryResponse)
		}

		c.JSON(http.StatusOK, gin.H{"data": categoryResponses})
	}
}

func CreateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categoryBody models.EditCategoryBody
		var category models.Category

		if err := c.ShouldBindJSON(&categoryBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category.SoldProductAmount = 0
		category.Type = categoryBody.Type

		_, err := govalidator.ValidateStruct(categoryBody)
		if err != nil {

			var errorMessages []string
			errs := err.(govalidator.Errors).Errors()

			for _, e := range errs {
				errorMessages = append(errorMessages, e.Error())
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
			return

		} else {
			categoryDB := db.Where("type = ?", &category.Type).First(&category)

			if categoryDB.Error == nil {
				c.JSON(http.StatusFound, gin.H{"error": "Category Already exist"})
				return
			}

			if result := db.Create(&category); result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}

			response := resCreateCategory{
				ID:                category.ID,
				Type:              category.Type,
				SoldProductAmount: category.SoldProductAmount,
				CreatedAt:         category.CreatedAt,
			}

			c.JSON(http.StatusCreated, gin.H{"data": response})
		}
	}
}

func EditCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category models.Category
		var editCategory models.EditCategoryBody
		categoryIdParam := c.Param("categoryId")

		categoryId, err := strconv.Atoi(categoryIdParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid param"})
			return
		}

		if err := c.ShouldBindJSON(&editCategory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.First(&category, categoryId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		db.Model(&category).Clauses(clause.Returning{}).Update("type", editCategory.Type)

		response := resCreateCategory{
			ID:                category.ID,
			Type:              category.Type,
			SoldProductAmount: category.SoldProductAmount,
			CreatedAt:         category.CreatedAt,
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func DeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category models.Category
		categoryIdParam := c.Param("categoryId")

		categoryId, err := strconv.Atoi(categoryIdParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid param"})
			return
		}

		if err := db.First(&category, categoryId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := db.Delete(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category has been successfully deleted"})
	}
}
