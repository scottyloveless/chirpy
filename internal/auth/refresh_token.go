package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)

	refreshToken := hex.EncodeToString(key)

	return refreshToken
}
