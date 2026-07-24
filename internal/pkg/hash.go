package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// HashPasswordFE hashes password using frontend format: salt:password
func HashPasswordFE(password, salt string) string {
	data := []byte(salt + ":" + password)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// HashPasswordLegacy hashes password using legacy format: password+salt
func HashPasswordLegacy(password, salt string) string {
	hHasher := sha256.New()
	hHasher.Write([]byte(password + salt))
	return hex.EncodeToString(hHasher.Sum(nil))
}

// GenerateSalt generates a random 16-byte hex salt
func GenerateSalt() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
