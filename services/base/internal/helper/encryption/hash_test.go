package encryption

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {
	password := "mySecurePassword123"
	
	hash, err := HashPassword(password)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
	// bcrypt hashes start with $2a$, $2b$, or $2y$
	assert.True(t, strings.HasPrefix(hash, "$2"))
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	password := ""
	
	hash, err := HashPassword(password)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestHashPassword_LongPassword(t *testing.T) {
	// Test with a 72-byte password (bcrypt's limit)
	password := strings.Repeat("a", 72)
	
	hash, err := HashPassword(password)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestHashPassword_SamePasswordDifferentHashes(t *testing.T) {
	password := "samePassword"
	
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	// Bcrypt generates different salts, so hashes should be different
	assert.NotEqual(t, hash1, hash2)
}

func TestComparePassword_CorrectPassword(t *testing.T) {
	password := "correctPassword"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	
	err = ComparePassword(hashedPassword, password)
	
	assert.NoError(t, err)
}

func TestComparePassword_IncorrectPassword(t *testing.T) {
	password := "correctPassword"
	wrongPassword := "wrongPassword"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	
	err = ComparePassword(hashedPassword, wrongPassword)
	
	assert.Error(t, err)
}

func TestComparePassword_EmptyPassword(t *testing.T) {
	password := "somePassword"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	
	err = ComparePassword(hashedPassword, "")
	
	assert.Error(t, err)
}

func TestComparePassword_InvalidHash(t *testing.T) {
	invalidHash := "not-a-valid-bcrypt-hash"
	password := "somePassword"
	
	err := ComparePassword(invalidHash, password)
	
	assert.Error(t, err)
}

func TestComparePassword_EmptyHash(t *testing.T) {
	err := ComparePassword("", "somePassword")
	
	assert.Error(t, err)
}
