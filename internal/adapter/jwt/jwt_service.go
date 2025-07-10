package jwt

import (
	"errors"
	"time"

	"github.com/AndreySmirnoffv/golang-auth-training/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateAccessToken(userID uint) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	ValidateAccessToken(token string) (*jwt.MapClaims, error)
	ValidateRefreshToken(token string) (*jwt.MapClaims, error)
	GenerateTokens(user *entity.User) (string, string, error)
	RefreshTokens(refreshToken string) (string, string, error)
}

type jwtService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) JWTService {
	return &jwtService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (j *jwtService) GenerateTokens(user *entity.User) (string, string, error) {
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(j.accessTTL).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	at, err := accessToken.SignedString(j.accessSecret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(j.refreshTTL).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	rt, err := refreshToken.SignedString(j.refreshSecret)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (j *jwtService) GenerateAccessToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(j.accessTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.accessSecret)
}

func (j *jwtService) GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(j.refreshTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.refreshSecret)
}

func (j *jwtService) ValidateAccessToken(tokenStr string) (*jwt.MapClaims, error) {
	return j.validateToken(tokenStr, j.accessSecret)
}

func (j *jwtService) ValidateRefreshToken(tokenStr string) (*jwt.MapClaims, error) {
	return j.validateToken(tokenStr, j.refreshSecret)
}

func (j *jwtService) validateToken(tokenStr string, secret []byte) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return &claims, nil
}

func (j *jwtService) RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := j.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	userIDFloat, ok := (*claims)["user_id"].(float64)
	if !ok {
		return "", "", errors.New("invalid refresh token claims")
	}

	userID := uint(userIDFloat)

	newAccess, err := j.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	newRefresh, err := j.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}
