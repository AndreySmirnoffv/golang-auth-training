package usecases

import (
	"errors"

	"github.com/AndreySmirnoffv/golang-auth-training/internal/adapter/jwt"
	"github.com/AndreySmirnoffv/golang-auth-training/internal/entity"
	"github.com/AndreySmirnoffv/golang-auth-training/pkg"
)

var (
	ErrEmailExists        = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type UserRepo interface {
	Create(u *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	Save(u *entity.User) error
	FindByID(id uint) (*entity.User, error)
}

type UserUseCase struct {
	repo UserRepo
	jwt  jwt.JWTService
}

func NewUserUseCase(r UserRepo, jwtSrv jwt.JWTService) *UserUseCase {
	return &UserUseCase{repo: r, jwt: jwtSrv}
}

func (uc *UserUseCase) Register(u *entity.User) error {
	if existing, _ := uc.repo.FindByEmail(u.Email); existing != nil {
		return ErrEmailExists
	}

	hashed, err := pkg.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashed

	return uc.repo.Create(u)
}

func (uc *UserUseCase) Login(email, password string) (*entity.User, string, string, error) {
	user, err := uc.repo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, "", "", ErrInvalidCredentials
	}

	if !pkg.CheckPasswordHash(password, user.Password) {
		return nil, "", "", ErrInvalidCredentials
	}

	accessToken, refreshToken, err := uc.jwt.GenerateTokens(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (uc *UserUseCase) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := uc.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	userIDFloat, ok := (*claims)["user_id"].(float64)
	if !ok {
		return "", "", ErrInvalidCredentials
	}
	userID := uint(userIDFloat)

	user, err := uc.repo.FindByID(userID)
	if err != nil || user == nil {
		return "", "", ErrUserNotFound
	}

	return uc.jwt.GenerateTokens(user)
}
