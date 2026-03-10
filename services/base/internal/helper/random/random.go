package random

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomCode returns a cryptographically random uppercase alphanumeric
// string of the given length. Uses crypto/rand — not seeded, never repeats.
func GenerateRandomCode(length int) string {
	code := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))
	for i := range code {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			panic("random: crypto/rand unavailable: " + err.Error())
		}
		code[i] = charset[n.Int64()]
	}
	return string(code)
}

func GenerateSecureToken(n int) (rawToken string, hashedToken string, err error) {
	b := make([]byte, n)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}

	rawToken = base64.URLEncoding.EncodeToString(b)
	hash := sha256.Sum256([]byte(rawToken))
	hashedToken = hex.EncodeToString(hash[:])
	return rawToken, hashedToken, nil
}
