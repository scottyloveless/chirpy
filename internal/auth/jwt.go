package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(60) * time.Minute)),
		Subject:   fmt.Sprintf("%v", userID),
	})

	sString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("Error signing string: %v", err)
		return "", err
	}

	return sString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) { return []byte(tokenSecret), nil })
	if err != nil {
		log.Printf("Token invalid or expired: %v", err)
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		log.Printf("Error fetching userID from claims: %v", err)
		return uuid.Nil, err
	}

	parsedID, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Printf("Error parsing userID: %v", err)
		return uuid.Nil, err
	}

	return parsedID, nil
}
