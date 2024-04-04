package models

import "github.com/dgrijalva/jwt-go"

// ---------------- DTO ----------------

// JWTClaim represents the authorized object encrypted in the JWT token
type JWTClaim struct {
	UserID    uint   `json:"user_id"`
	UserEmail string `json:"user_email"`
	jwt.StandardClaims
}

type RegisterRequest struct {
	Name              string `json:"name" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	Password          string `json:"password" validate:"required"`
	BirthdayInSeconds int64  `json:"birthday_in_seconds" validate:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
