package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		log.Printf("No authorization header present")
		return "", errors.New("no authorization header present")
	}

	if !strings.HasPrefix(auth, "ApiKey") {
		log.Printf("Malformed Authorization header")
		return "", errors.New("malformed authorization header")
	}

	return strings.TrimSpace(strings.TrimPrefix(auth, "ApiKey")), nil
}
