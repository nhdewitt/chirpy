package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		log.Printf("No Authorization header present")
		return "", errors.New("authorization header doesn't exist")
	}

	if !strings.HasPrefix(auth, "Bearer") {
		log.Printf("Malformed Authorization header")
		return "", errors.New("malformed authorization header")
	}

	return strings.TrimSpace(strings.TrimPrefix(auth, "Bearer")), nil
}
