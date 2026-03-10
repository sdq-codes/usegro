package random

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomCode_CorrectLength(t *testing.T) {
	lengths := []int{4, 6, 8, 10, 16}
	
	for _, length := range lengths {
		t.Run("", func(t *testing.T) {
			code := GenerateRandomCode(length)
			assert.Equal(t, length, len(code))
		})
	}
}

func TestGenerateRandomCode_OnlyValidCharacters(t *testing.T) {
	code := GenerateRandomCode(100)
	
	for _, char := range code {
		assert.Contains(t, charset, string(char))
	}
}

func TestGenerateRandomCode_ZeroLength(t *testing.T) {
	code := GenerateRandomCode(0)
	assert.Equal(t, 0, len(code))
	assert.Equal(t, "", code)
}

func TestGenerateRandomCode_Uniqueness(t *testing.T) {
	// Generate multiple codes and ensure they're not all the same
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code := GenerateRandomCode(10)
		codes[code] = true
	}
	
	// With 100 codes of length 10, we should have multiple unique values
	// (collision is possible but extremely unlikely)
	assert.Greater(t, len(codes), 90)
}

func TestGenerateSecureToken_Success(t *testing.T) {
	n := 32
	
	rawToken, hashedToken, err := GenerateSecureToken(n)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, rawToken)
	assert.NotEmpty(t, hashedToken)
}

func TestGenerateSecureToken_RawTokenIsBase64(t *testing.T) {
	rawToken, _, err := GenerateSecureToken(32)
	
	assert.NoError(t, err)
	
	// Should be decodable as base64
	decoded, err := base64.URLEncoding.DecodeString(rawToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, decoded)
}

func TestGenerateSecureToken_HashedTokenIsHex(t *testing.T) {
	_, hashedToken, err := GenerateSecureToken(32)
	
	assert.NoError(t, err)
	
	// Should be valid hex
	decoded, err := hex.DecodeString(hashedToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, decoded)
	// SHA256 produces 32 bytes (64 hex characters)
	assert.Equal(t, 64, len(hashedToken))
}

func TestGenerateSecureToken_Uniqueness(t *testing.T) {
	// Generate multiple tokens and ensure they're unique
	tokens := make(map[string]bool)
	hashes := make(map[string]bool)
	
	for i := 0; i < 100; i++ {
		rawToken, hashedToken, err := GenerateSecureToken(32)
		assert.NoError(t, err)
		tokens[rawToken] = true
		hashes[hashedToken] = true
	}
	
	// All tokens and hashes should be unique
	assert.Equal(t, 100, len(tokens))
	assert.Equal(t, 100, len(hashes))
}

func TestGenerateSecureToken_DifferentSizes(t *testing.T) {
	sizes := []int{16, 32, 64, 128}
	
	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			rawToken, hashedToken, err := GenerateSecureToken(size)
			
			assert.NoError(t, err)
			assert.NotEmpty(t, rawToken)
			assert.NotEmpty(t, hashedToken)
			// Hash should always be 64 characters (SHA256)
			assert.Equal(t, 64, len(hashedToken))
		})
	}
}

func TestGenerateSecureToken_ZeroSize(t *testing.T) {
	rawToken, hashedToken, err := GenerateSecureToken(0)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedToken)
	// Even with 0 bytes, we still get a hash of the empty string
	assert.Equal(t, 64, len(hashedToken))
}
