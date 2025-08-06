package auth_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/nhdewitt/chirpy/internal/auth"
	"github.com/stretchr/testify/require"
)

func TestGetBearerTokenBearerExists(t *testing.T) {
	userID := uuid.New()
	// Generate token
	token, err := auth.MakeJWT(userID, os.Getenv("CHIRPY_AUTH_SECRET"), 30)
	require.NoError(t, err, "MakeJWT should not return an error")
	header := make(http.Header)
	header.Set("Authorization", "Bearer    "+token)
	require.NotEmpty(t, header, "Header should not be empty")

	// Validate token
	bearer, err := auth.GetBearerToken(header)
	require.NoError(t, err, "GetBearerToken should not return an error")
	require.Equal(t, token, bearer)
}

func TestGetBearerTokenNoHeaderRaisesError(t *testing.T) {
	header := make(http.Header)

	_, err := auth.GetBearerToken(header)
	require.Error(t, err, "Missing header should return an error")
}

func TestGetBearerTokenMalformedHeaderRaisesError(t *testing.T) {
	userID := uuid.New()
	header := make(http.Header)

	token, err := auth.MakeJWT(userID, os.Getenv("CHIRPY_AUTH_SECRET"), 30)
	require.NoError(t, err, "MakeJWT should not return an error")

	header.Set("Authorization", "earer    "+token)
	_, err = auth.GetBearerToken(header)
	require.Error(t, err, "Malformed header should return an error")
}
