package dto

type EmailVerificationDTI struct {
	TokenHash string `json:"token_hash" validate:"required,min=6,max=6"`
}
