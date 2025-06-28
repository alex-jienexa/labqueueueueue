package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var secretKey []byte

func GetSecretKey() []byte {
	if secretKey == nil {
		godotenv.Load()
		if os.Getenv("JWT_SECRET") == "" {
			secretKey = []byte("not-so-secret-key")
		}
		secretKey = []byte(os.Getenv("JWT_SECRET"))
		return secretKey
	} else {
		return secretKey
	}
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetSecretKey())
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*Claims), nil
}
