package models

import "github.com/dgrijalva/jwt-go"

// ---------------- DTO ----------------

// JWTClaim represents the authorized object encrypted in the JWT token
type JWTClaim struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	jwt.StandardClaims
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyEmailRequest struct {
	Email            string `json:"email" validate:"required,email"`
	VerificationCode string `json:"verification_code" validate:"required"`
}

type SendVerificationEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}
