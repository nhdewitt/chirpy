package auth_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/nhdewitt/chirpy/internal/auth"
)

const (
	secret      = "supersecret"
	wrongSecret = "wrongsecret"
	emptySecret = ""
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	expiresIn := time.Minute

	// Generate token
	token, err := auth.MakeJWT(userID, secret, expiresIn)
	require.NoError(t, err, "MakeJWT should not return an error")
	require.NotEmpty(t, token, "Token should not be empty")

	// Validate token
	parsedID, err := auth.ValidateJWT(token, secret)
	require.NoError(t, err, "ValidateJWT should not return an error")
	require.Equal(t, userID, parsedID, "Parsed userID should match original userID")
}

func TestValidateJWT_InvalidClaimsType(t *testing.T) {
	// Token with map claims instead of RegisteredClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "someid",
	})
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	_, err = auth.ValidateJWT(tokenString, secret)
	require.Error(t, err, "ValidateJWT should fail when claims type is not RegisteredClaims")
}

func TestValidateJWT_InvalidUUIDSubject(t *testing.T) {
	// Token with invalid UUID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: "not-a-uuid",
	})
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	_, err = auth.ValidateJWT(tokenString, secret)
	require.Error(t, err, "ValidateJWT should with with invalid UUID")
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()

	// Create expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
	})
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	_, err = auth.ValidateJWT(tokenString, secret)
	require.Error(t, err, "ValidateJWT should fail for expired token")
}

func TestMakeJWT_SignError(t *testing.T) {
	userID := uuid.New()

	// Override signing method to invalidate
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
		Subject: userID.String(),
	})
	_, err := token.SignedString([]byte(emptySecret))
	require.Error(t, err, "Signing with invalid method should error")
}

func TestValidateJWT_InvalidSecret(t *testing.T) {
	userID := uuid.New()

	// Create token with correct secret
	token, err := auth.MakeJWT(userID, secret, time.Minute)
	require.NoError(t, err)

	// Validate with wrong secret
	_, err = auth.ValidateJWT(token, wrongSecret)
	require.Error(t, err, "ValidateJWT should fail with wrong secret")
}

func TestMakeJWT_Expiration(t *testing.T) {
	userID := uuid.New()

	// Token expires immediately
	token, err := auth.MakeJWT(userID, secret, 0)
	require.NoError(t, err)

	// Parse token to inspect claims
	parsedToken, _ := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	require.True(t, ok)
	require.WithinDuration(t, time.Now().UTC(), claims.IssuedAt.Time, time.Second)
	require.WithinDuration(t, time.Now().UTC(), claims.ExpiresAt.Time, time.Second)
}
