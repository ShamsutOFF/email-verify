package hash

import (
	"crypto/rand"
	"encoding/hex"
)

// Generate создаёт случайную 16-байтную строку в HEX-формате
func Generate() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
