package dto

import "github.com/google/uuid"

type RegisterUserDTI struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64,passwd_lower,passwd_upper,passwd_digit,passwd_special"`
}

type UserExistDTI struct {
	Email string `json:"email" validate:"required,email"`
}

type RegisteredUserDTO struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type ForgotPasswordDTI struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordDTI struct {
	Token           string `json:"token" validate:"required"`
	Password        string `json:"password" validate:"required,min=8,max=64,passwd_lower,passwd_upper,passwd_digit,passwd_special"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type RefreshTokenDTI struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthTokensDTO struct {
	ID           interface{} `json:"id"`
	Email        string      `json:"email"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

// GoogleUserInfo is the response from https://www.googleapis.com/oauth2/v2/userinfo
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type RequestEmailCodeDTI struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyEmailCodeDTI struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}
