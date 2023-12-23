package handlers

import (
	"admin-dashboard-FP/models"
	"errors"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type resPostTransaction struct {
	TotalPrice   int    `json:"total_price"`
	Quantity     int    `json:"quantity"`
	ProductTitle string `json:"product_title"`
}

type resGetAllTransactionUser struct {
	ID         uint `json:"id"`
	ProductID  uint `json:"product_id"`
	UserID     uint `json:"user_id"`
	Quantity   int  `json:"quantity"`
	TotalPrice int  `json:"total_price"`
	Product    GetProductTransaction
}

type resGetAllTransactionAdmin struct {
	ID         uint `json:"id"`
	ProductID  uint `json:"product_id"`
	UserID     uint `json:"user_id"`
	Quantity   int  `json:"quantity"`
	TotalPrice int  `json:"total_price"`
	Product    GetProductTransaction
	User       GetUserTransaction
}

type GetProductTransaction struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Price      int       `json:"price"`
	Stock      int       `json:"stock"`
	CategoryId uint      `json:"category_Id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GetUserTransaction struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetTransactionLoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactionHistory []models.TransactionHistory

		dataToken, _ := c.Get("Token")
		userID := dataToken.(*models.DataToken).ID

		if err := db.Preload("Product").Where("user_id = ?", userID).Find(&transactionHistory).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Transaction History User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var transactionResponses []resGetAllTransactionUser

		for _, eachTransaction := range transactionHistory {
			var productResponse GetProductTransaction
			var transactionResponse resGetAllTransactionUser

			productResponse = GetProductTransaction{
				ID:         eachTransaction.Product.ID,
				Title:      eachTransaction.Product.Title,
				Price:      eachTransaction.Product.Price,
				Stock:      eachTransaction.Product.Stock,
				CategoryId: eachTransaction.Product.CategoryID,
				CreatedAt:  eachTransaction.Product.CreatedAt,
				UpdatedAt:  eachTransaction.Product.UpdatedAt,
			}

			transactionResponse = resGetAllTransactionUser{
				ID:         eachTransaction.ID,
				ProductID:  eachTransaction.ProductID,
				UserID:     eachTransaction.UserID,
				Quantity:   eachTransaction.Quantity,
				TotalPrice: eachTransaction.TotalPrice,
				Product:    productResponse,
			}

			transactionResponses = append(transactionResponses, transactionResponse)
		}

		c.JSON(http.StatusOK, gin.H{"Transactions": transactionResponses})
	}
}

func GetTransactionUserByAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transactionHistory []models.TransactionHistory

		if err := db.Preload("Product").Preload("User").Find(&transactionHistory).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Transaction History User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var transactionResponses []resGetAllTransactionAdmin

		for _, eachTransaction := range transactionHistory {
			var productResponse GetProductTransaction
			var transactionResponse resGetAllTransactionAdmin
			var userResponse GetUserTransaction

			productResponse = GetProductTransaction{
				ID:         eachTransaction.Product.ID,
				Title:      eachTransaction.Product.Title,
				Price:      eachTransaction.Product.Price,
				Stock:      eachTransaction.Product.Stock,
				CategoryId: eachTransaction.Product.CategoryID,
				CreatedAt:  eachTransaction.Product.CreatedAt,
				UpdatedAt:  eachTransaction.Product.UpdatedAt,
			}

			userResponse = GetUserTransaction{
				ID:        eachTransaction.User.ID,
				Email:     eachTransaction.User.Email,
				FullName:  eachTransaction.User.FullName,
				Balance:   eachTransaction.User.Balance,
				CreatedAt: eachTransaction.User.CreatedAt,
				UpdatedAt: eachTransaction.User.UpdatedAt,
			}

			transactionResponse = resGetAllTransactionAdmin{
				ID:         eachTransaction.ID,
				ProductID:  eachTransaction.ProductID,
				UserID:     eachTransaction.UserID,
				Quantity:   eachTransaction.Quantity,
				TotalPrice: eachTransaction.TotalPrice,
				Product:    productResponse,
				User:       userResponse,
			}

			transactionResponses = append(transactionResponses, transactionResponse)
		}

		c.JSON(http.StatusOK, gin.H{"Transactions_Users": transactionResponses})
	}
}

func PostTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyTransaction models.PostTransactionBody
		var transaction models.TransactionHistory
		var product models.Product
		var user models.User
		var category models.Category

		dataToken, _ := c.Get("Token")
		userID := dataToken.(*models.DataToken).ID

		if err := c.ShouldBindJSON(&bodyTransaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			_, err := govalidator.ValidateStruct(bodyTransaction)
			if err != nil {

				var errorMessages []string
				errs := err.(govalidator.Errors).Errors()

				for _, e := range errs {
					errorMessages = append(errorMessages, e.Error())
				}

				c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
				return

			} else {
				if err := db.First(&product, bodyTransaction.ProductId).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				if bodyTransaction.Quantity > product.Stock {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock available"})
					return
				}

				if err := db.First(&user, userID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				if bodyTransaction.Quantity*product.Price > user.Balance {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough balance, please top up first"})
					return
				}

				if err := db.First(&category, product.CategoryID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				remainingStock := product.Stock - bodyTransaction.Quantity
				remainingBalance := user.Balance - (bodyTransaction.Quantity * product.Price)
				alreadySoldCategory := category.SoldProductAmount + bodyTransaction.Quantity

				product.Stock = remainingStock
				user.Balance = remainingBalance
				category.SoldProductAmount = alreadySoldCategory

				if err := db.Save(&product).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				if err := db.Save(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				if err := db.Save(&category).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				transaction.ProductID = bodyTransaction.ProductId
				transaction.UserID = userID
				transaction.Quantity = bodyTransaction.Quantity
				transaction.TotalPrice = (bodyTransaction.Quantity * product.Price)

				if err := db.Save(&transaction).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				response := resPostTransaction{
					TotalPrice:   transaction.TotalPrice,
					Quantity:     bodyTransaction.Quantity,
					ProductTitle: product.Title,
				}

				c.JSON(http.StatusCreated, gin.H{
					"message":          "You have succesfully purchased the product",
					"transaction_bill": response,
				})
			}
		}
	}
}
