package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/daodao97/egin/egin/utils/config"
)

var jwtSecret = []byte(config.Config.Jwt.Secret)

type Claims struct {
	UserID    int    `json:"user_id"`
	UserEmail string `json:"user_email"`
	jwt.StandardClaims
}

func GenerateToken(id int, email string) (string, error) {
	claims := Claims{
		id,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + config.Config.Jwt.TokenExpire,
			Issuer:    "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
