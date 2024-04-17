package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/repositories"
	"github.com/hodukihugi/winglets-api/utils"
	"time"
)

type IAuthService interface {
	Authorize(string) (*models.JWTClaim, error)
	GenerateJWTTokens(user models.User) (string, string, int64, int64, error)
	Register(request models.RegisterRequest) error
	Refresh(user models.User) (string, int64, error)
}

// AuthService service relating to authorization
type AuthService struct {
	env      *core.Env
	logger   *core.Logger
	userRepo repositories.IUserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(
	env *core.Env,
	logger *core.Logger,
	userRepo repositories.IUserRepository,
) IAuthService {
	return &AuthService{
		env:      env,
		logger:   logger,
		userRepo: userRepo,
	}
}

// Authorize authorizes the generated token
func (s *AuthService) Authorize(tokenString string) (*models.JWTClaim, error) {
	var claim jwt.Claims = &models.JWTClaim{}

	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.env.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claim.(*models.JWTClaim), nil
}

func (s *AuthService) Refresh(user models.User) (string, int64, error) {
	return s.signJWT(user, s.env.AccessTokenExpiresIn)
}

// GenerateJWTTokens creates jwt auth tokens
func (s *AuthService) GenerateJWTTokens(user models.User) (string, string, int64, int64, error) {
	accessToken, accessExpired, err := s.signJWT(user, s.env.AccessTokenExpiresIn)
	if err != nil {
		return "", "", 0, 0, err
	}

	refreshToken, refreshExpired, err := s.signJWT(user, s.env.RefreshTokenExpiresIn)
	if err != nil {
		return "", "", 0, 0, err
	}

	return accessToken, refreshToken, accessExpired, refreshExpired, nil
}

func (s *AuthService) Register(request models.RegisterRequest) error {
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}

	return s.userRepo.Create(models.User{
		Email:    request.Email,
		Password: hashedPassword,
	})
}

// ----------------- private -----------------

func (s *AuthService) signJWT(user models.User, expirationPeriod time.Duration) (string, int64, error) {
	now := time.Now().UTC()
	exp := now.Add(expirationPeriod).Unix()
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTClaim{
		UserID:    user.ID,
		UserEmail: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Issuer:    "winglets-web",
		},
	}).SignedString([]byte(s.env.JWTSecret))
	return jwtToken, exp, err
}
