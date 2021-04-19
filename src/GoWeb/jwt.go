package main

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const JWT_TOKEN_VALIDITY = 24 * 60 * 60 * 1000

type CustomClaims struct {
	Subject   string
	Id        int
	ExpiresAt int64
	IssuedAt  int64
}

func (c CustomClaims) Valid() error {
	now := time.Now().UnixNano()

	//密钥过期
	if c.ExpiresAt < now/1e6 {
		return JwtError{message: "JWT token has expired"}
	}

	return nil
}

func generateJwtToken(user AuthUser) string {
	currentTime := time.Now().UnixNano() / 1e6
	claims := CustomClaims{
		Subject:   user.Username,
		Id:        user.Id,
		ExpiresAt: currentTime + JWT_TOKEN_VALIDITY,
		IssuedAt:  currentTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, _ := token.SignedString([]byte("min107"))
	return "Bearer " + tokenString
}

func validateJwtToken(tokenString string, claims *CustomClaims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("min107"), nil
	})
	if err != nil {
		return err
	}

	if token != nil && token.Valid {
		return err
	}
	return nil

}
