package random

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

// NOTE: use this as the generator of all object ids.
func RandID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	id := hex.EncodeToString(b)
	return id, nil
}

func RandBase64(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
