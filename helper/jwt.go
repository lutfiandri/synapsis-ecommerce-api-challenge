package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
)

var JWT_KEY = []byte(os.Getenv("JWT_SECRET_KEY"))

type JWTClaim struct {
	UserID    uint
	UserEmail string
	UserRole  string
	jwt.RegisteredClaims
}

func GenerateJWT(user *model.User) (string, error) {
	now := time.Now()
	expTime := now.Add(time.Hour * 12)

	claims := &JWTClaim{
		UserID:    user.ID,
		UserEmail: user.Email,
		UserRole:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "synapsis-ecommerce-api-challenge",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenAlgo.SignedString(JWT_KEY)

	return token, err
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	claims := &JWTClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})

	return token, err
}
