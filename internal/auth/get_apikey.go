package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("missing Authorization header")
	}

	if !strings.HasPrefix(tokenString, "ApiKey ") {
		return "", errors.New("invalid Authorization header format")
	}

	stripped := strings.TrimSpace(strings.TrimPrefix(tokenString, "ApiKey "))
	if stripped == "" {
		return "", errors.New("empty api key")
	}
	return stripped, nil
}
