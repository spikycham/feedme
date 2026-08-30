package random

import (
	"crypto/rand"
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
