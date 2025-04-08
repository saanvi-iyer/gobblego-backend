package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/saanvi-iyer/gobblego-backend/models"
)

func GenerateToken(user models.User) (string, error) {

	claims := jwt.MapClaims{
		"user_id":   user.UserID.String(),
		"cart_id":   user.CartID.String(),
		"user_name": user.UserName,
		"is_leader": user.IsLeader,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// TODO: Use environment variable for secret key in production
	secretKey := []byte("abcdef")
	tokenString, err := token.SignedString(secretKey)

	return tokenString, err
}
