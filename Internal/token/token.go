// Package token: работа с токеном jwt
package token

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"Service_1Cv8/internal/constants"
)

type ClaimStore struct {
	Key    string
	Value  []byte
	Secret []byte
}

type ClaimsStore map[string]ClaimStore

// Claim данные токена, полусенные с сервера
type Claim struct {
	Authorized bool
	Key        string
	Exp        float64
}

// GenerateJWT генерация токена для пользователя
func (c *Claim) GenerateJWT(s []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = c.Authorized
	claims["key"] = c.Key
	claims["exp"] = c.Exp

	if len(s) == 0 {
		s = []byte(constants.HashKey)
	}

	tokenString, err := token.SignedString(s)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

// ExtractClaims получение имя пользователя из токена
func ExtractClaims(tokenStr string, secret []byte) (jwt.MapClaims, bool) {
	//hmacSecret := []byte(constants.HashKey)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Println("Invalid JWT Token")
		return nil, false
	}
}

// NewClaims Инициализация сущности Claims
func NewClaims(key string, timeLiveToken time.Duration) *Claim {
	return &Claim{
		Authorized: true,
		Key:        key,
		Exp:        float64(time.Now().Add(time.Hour * 24 * timeLiveToken).Unix()),
	}
}
