package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

func GetDBConfig() *DBConfig {
	if os.Getenv("ENVIROMENT") == "local" {
		return &DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME_FP"),
		}
	} else {
		return &DBConfig{
			Host:     os.Getenv("DB_HOST_PROD"),
			Port:     os.Getenv("DB_PORT_PROD"),
			User:     os.Getenv("DB_USER_PROD"),
			Password: os.Getenv("DB_PASSWORD_PROD"),
			DBName:   os.Getenv("DB_NAME_FP_PROD"),
		}
	}
}

func (config *DBConfig) GetDBURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		config.Host, config.Port, config.User, config.DBName, config.Password)
}
