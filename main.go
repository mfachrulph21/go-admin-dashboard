package main

import (
	"admin-dashboard-FP/config"
	"admin-dashboard-FP/models"
	"admin-dashboard-FP/routes"
	"log"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//config db connection
	dbConfig := config.GetDBConfig()
	dsn := dbConfig.GetDBURL()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database :", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{}, &models.TransactionHistory{}); err != nil {
		log.Fatal("Failed to automigrate :", err)
	}

	//add custom validator
	govalidator.TagMap["role"] = govalidator.Validator(func(str string) bool {
		return str == "admin" || str == "customer"
	})

	govalidator.TagMap["minimalStock"] = govalidator.Validator(func(str string) bool {
		num, _ := strconv.Atoi(str)
		return num >= 5
	})

	router := gin.Default()

	routes.UserRoutes(router, db)
	routes.ProductRoutes(router, db)
	routes.CategoryRoutes(router, db)
	routes.TransactionRoutes(router, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Run server failed : ", err)
	}
}
