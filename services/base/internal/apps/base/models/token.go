package models

import jwt "github.com/golang-jwt/jwt/v5"

// TokenClaims is the JWT payload — contains only what is needed for auth,
// not the full User struct (avoids leaking sensitive fields into every token).
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
