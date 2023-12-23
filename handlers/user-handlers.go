package handlers

import (
	"admin-dashboard-FP/helpers"
	"admin-dashboard-FP/models"
	"fmt"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type resRegister struct {
	ID        uint      `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func UserRegister(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser models.User

		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			newUser.Role = "customer"
			newUser.Balance = 0

			_, err := govalidator.ValidateStruct(newUser)
			if err != nil {

				var errorMessages []string
				errs := err.(govalidator.Errors).Errors()

				for _, e := range errs {
					errorMessages = append(errorMessages, e.Error())
				}

				c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
				return
			} else {
				userDB := db.Where("email = ?", &newUser.Email).First(&newUser)

				if userDB.Error == nil {
					c.JSON(http.StatusFound, gin.H{"error": "User Already exist"})
					return
				}

				stringValue, err := helpers.HashPassword(newUser.Password)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				newUser.Password = stringValue

				if result := db.Create(&newUser); result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
					return
				}

				response := resRegister{
					ID:        newUser.ID,
					FullName:  newUser.FullName,
					Email:     newUser.Email,
					Password:  newUser.Password,
					Balance:   newUser.Balance,
					CreatedAt: newUser.CreatedAt,
				}

				c.JSON(http.StatusCreated, gin.H{"data": response})
			}
		}
	}
}

func UserLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUser models.LoginUser
		var users models.User

		if err := c.ShouldBindJSON(&loginUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := govalidator.ValidateStruct(loginUser)
		if err != nil {

			var errorMessages []string
			errs := err.(govalidator.Errors).Errors()

			for _, e := range errs {
				errorMessages = append(errorMessages, e.Error())
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
			return
		}

		result := db.Where("email = ?", loginUser.Email).Find(&users)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "wrong email/password"})
			return
		} else {
			match := helpers.CheckPasswordHash(loginUser.Password, users.Password)

			if !match {
				c.JSON(http.StatusNotFound, gin.H{"error": "wrong email/password"})
				return
			}

			token, err := helpers.GenerateJWT(users.ID, users.Role, loginUser)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": token})
		}
	}
}

func UserTopup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var topupBody models.TopupBody
		dataToken, _ := c.Get("Token")
		userID := dataToken.(*models.DataToken).ID

		if err := c.ShouldBindJSON(&topupBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := govalidator.ValidateStruct(topupBody)
		if err != nil {

			var errorMessages []string
			errs := err.(govalidator.Errors).Errors()

			for _, e := range errs {
				errorMessages = append(errorMessages, e.Error())
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
			return
		}

		result := db.First(&user, userID)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		updatedBalance := user.Balance + topupBody.Balance

		db.Model(&user).Update("balance", updatedBalance)

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Your balance has been updated to Rp. %d", updatedBalance)})
	}
}
