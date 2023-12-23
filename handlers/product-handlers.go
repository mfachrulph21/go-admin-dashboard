package handlers

import (
	"admin-dashboard-FP/helpers"
	"admin-dashboard-FP/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type resCreateProduct struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Price      int       `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID uint      `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type resEditProduct struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Price      any       `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID uint      `json:"CategoryId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func GetProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.Product

		if err := db.Find(&products).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "no product found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []resCreateProduct

		for _, product := range products {
			item := resCreateProduct{
				ID:         product.ID,
				Title:      product.Title,
				Price:      product.Price,
				Stock:      product.Stock,
				CategoryID: product.CategoryID,
				CreatedAt:  product.CreatedAt,
			}

			response = append(response, item)
		}

		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProduct models.Product
		var category models.Category

		if err := c.ShouldBindJSON(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			_, err := govalidator.ValidateStruct(newProduct)
			if err != nil {

				var errorMessages []string
				errs := err.(govalidator.Errors).Errors()

				for _, e := range errs {
					errorMessages = append(errorMessages, e.Error())
				}

				c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
				return

			} else {
				categoryDB := db.Where("id = ?", &newProduct.CategoryID).First(&category)

				if categoryDB.Error != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Category ID doesn't exist"})
					return
				}

				if result := db.Create(&newProduct); result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
					return
				}

				response := resCreateProduct{
					ID:         newProduct.ID,
					Title:      newProduct.Title,
					Price:      newProduct.Price,
					Stock:      newProduct.Stock,
					CategoryID: newProduct.CategoryID,
					CreatedAt:  newProduct.CreatedAt,
				}

				c.JSON(http.StatusCreated, gin.H{"data": response})
			}
		}
	}
}

func EditProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updatedProduct models.EditProductBody
		var product models.Product
		var category models.Category

		productIdParam := c.Param("productId")

		productId, err := strconv.Atoi(productIdParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid param",
			})
			return
		}

		if err := c.ShouldBindJSON(&updatedProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			_, err := govalidator.ValidateStruct(updatedProduct)
			if err != nil {

				var errorMessages []string
				errs := err.(govalidator.Errors).Errors()

				for _, e := range errs {
					errorMessages = append(errorMessages, e.Error())
				}

				c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
				return
			} else {
				if err := db.First(&category, updatedProduct.CategoryID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						c.JSON(http.StatusNotFound, gin.H{"error": "Category ID not found"})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				if err := db.First(&product, productId).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				product.Title = updatedProduct.Title
				product.Price = updatedProduct.Price
				product.Stock = updatedProduct.Stock
				product.CategoryID = updatedProduct.CategoryID

				if err := db.Save(&product).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				response := resEditProduct{
					ID:         product.ID,
					Title:      product.Title,
					Price:      helpers.FormatUang(product.Price),
					Stock:      product.Stock,
					CategoryID: product.CategoryID,
					CreatedAt:  product.CreatedAt,
					UpdatedAt:  product.UpdatedAt,
				}

				c.JSON(http.StatusOK, gin.H{"product": response})
			}
		}
	}
}

func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		productIdParam := c.Param("productId")

		productId, err := strconv.Atoi(productIdParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid param"})
			return
		}

		if err := db.First(&product, productId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := db.Delete(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product has been successfully deleted"})
	}
}
