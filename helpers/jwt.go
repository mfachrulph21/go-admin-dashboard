package helpers

import (
	"admin-dashboard-FP/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userID uint, userRole string, userInput models.LoginUser) (string, error) {
	var secretKey string = os.Getenv("SecretKey")
	expirationTime := time.Now().Add(time.Hour * 1)

	dataToken := &models.DataToken{
		ID:    userID,
		Email: userInput.Email,
		Role:  userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, dataToken)
	signedToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
