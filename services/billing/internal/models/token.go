package models

import jwt "github.com/golang-jwt/jwt/v5"

type Token struct {
	User User `json:"user" bson:"user"`
	jwt.RegisteredClaims
}
