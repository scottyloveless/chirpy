package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJWT(t *testing.T) {
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	tokenSecret1 := "Testing123!"
	jwt1, err := MakeJWT(uuid1, tokenSecret1)
	if err != nil {
		t.Errorf("Something is wrong making jwt: %v", err)
		return
	}
	jwt2, err2 := MakeJWT(uuid2, "WrongPass1!")
	if err2 != nil {
		t.Errorf("Something is wrong making jwt: %v", err2)
		return
	}

	tests := []struct {
		name        string
		tokenSecret string
		jwtTest     string
		wantErr     bool
	}{
		{
			name:        "valid password",
			tokenSecret: tokenSecret1,
			jwtTest:     jwt1,
			wantErr:     false,
		},
		{
			name:        "different uuid",
			tokenSecret: tokenSecret1,
			jwtTest:     jwt2,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateJWT(tt.jwtTest, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
