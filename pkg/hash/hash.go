package hash

import (
	"crypto/rand"
	"encoding/hex"
)

// Generate создаёт случайную 16-байтную строку в HEX-формате
func Generate() string {
	bytes := make([]byte, 16) // 128 бит
	if _, err := rand.Read(bytes); err != nil {
		panic(err) // в реальном коде лучше вернуть ошибку
	}
	return hex.EncodeToString(bytes)
}
