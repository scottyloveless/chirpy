package auth

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectedErr   error
	}{
		{
			name:          "Valid Bearer token",
			authHeader:    "Bearer abc123",
			expectedToken: "abc123",
			expectedErr:   nil,
		},
		{
			name:        "Missing Authorization header",
			authHeader:  "",
			expectedErr: errors.New("missing Authorization header"),
		},
		{
			name:        "Invalid header format",
			authHeader:  "Basic abc123",
			expectedErr: errors.New("invalid Authorization header format"),
		},
		{
			name:        "Empty token after Bearer",
			authHeader:  "Bearer ",
			expectedErr: errors.New("empty bearer token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			token, err := GetBearerToken(req.Header)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token != tt.expectedToken {
					t.Errorf("expected token %q, got %q", tt.expectedToken, token)
				}
			}
		})
	}
}
