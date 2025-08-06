package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func MakeRefreshToken() (string, error) {
	tokenString := make([]byte, 32)
	_, err := rand.Read(tokenString)
	if err != nil {
		return "", errors.New("error generating random data")
	}
	return hex.EncodeToString(tokenString), nil
}
