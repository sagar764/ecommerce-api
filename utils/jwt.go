package utils

import (
	"ecommerce-api/config"
	"ecommerce-api/internal/consts"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	secretKey []byte
)

func init() {
	cfg, err := config.LoadConfig(consts.AppName)
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}
	secretKey = []byte(cfg.JWTSecretKey)
}

// GenerateToken generates a JWT token
func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateToken parses and validates the JWT token
func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})
}
