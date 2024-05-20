// Package midware: middleware сервера. Проверка на аутетификацию
package midware

import (
	"Service_1Cv8/internal/constants"
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

// IsAuthorized middleware проверки пользователя по токену.
// Если токен валиден, то работа с данными разрешена
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "close")

		if r.Header["Authorization"] != nil {

			TokenFindMatches(endpoint, w, r)
			return
		}
		TokenNotFound(w)
	})
}

// TokenFindMatches проверки пользователя по токену.
func TokenFindMatches(endpoint func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {

	token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("there was an error")
		}
		return constants.HashKey, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		return
	}

	if token.Valid {
		endpoint(w, r)
	}
}

// TokenNotFound действие если токен не валиден
func TokenNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte("Not Authorized"))
	if err != nil {
		log.Println(err)
	}
}
