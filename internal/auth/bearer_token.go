package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	stripped := strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	if stripped == "" {
		return "", errors.New("empty bearer token")
	}
	return stripped, nil
}
