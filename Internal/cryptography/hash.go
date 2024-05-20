// Package cryptography: Хеширование паролей пользователей
package cryptography

import (
	"Service_1Cv8/internal/constants"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// HashSHA256 хеширует строку. Использует ключ. Ключ по умолчанию constants.HashKey
func HashSHA256(data string, strKey string) (hash string) {

	if strKey == "" {
		strKey = constants.HashKey
	}

	key := []byte(strKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	hash = fmt.Sprintf("%x", h.Sum(nil))
	return

}

func DecodeHashSHA256(data string) (hash string) {

	//key := []byte(strKey)
	//h := hmac.New(sha256.New, key)

	str, _ := hex.DecodeString(data)
	return string(str)

}
