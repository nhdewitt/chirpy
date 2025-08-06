package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	timeNow := jwt.NewNumericDate(time.Now().UTC())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  timeNow,
		ExpiresAt: jwt.NewNumericDate(timeNow.Add(expiresIn)),
		Subject:   userID.String(),
	})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("Can't sign token with key: %s", err)
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("Error parsing token: %s", err)
		return uuid.UUID{}, err
	}

	tok, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		log.Printf("Error getting ID: %s", err)
		return uuid.UUID{}, err
	}

	validated, err := uuid.Parse(tok.Subject)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		return uuid.UUID{}, err
	}

	return validated, nil
}
