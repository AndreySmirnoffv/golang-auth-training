package usecases

import (
	"errors"

	"github.com/AndreySmirnoffv/golang-auth-training/internal/entity"
	"github.com/AndreySmirnoffv/golang-auth-training/pkg"
)

var ErrEmailExists = errors.New("email already registered")
var ErrInavlidCredentials = errors.New("invalid credentials")

type UserRepo interface {
	Save(u *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase(r UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (uc *UserUseCase) Register(u *entity.User) error {
	if existing, _ := uc.repo.FindByEmail(u.Email); existing != nil {
		return ErrEmailExists
	}

	hashedPassword, err := pkg.HashPassword(u.Password)

	if err != nil {
		return err
	}

	u.Password = hashedPassword

	return uc.repo.Save(u)
}

func (uc *UserUseCase) Login(u *entity.User) error {
	existing, err := uc.repo.FindByEmail(u.Email)

	if err != nil {
		return err
	}

	if existing == nil {
		return ErrInavlidCredentials
	}

	ok := pkg.CheckPasswordHash(u.Password, existing.Password)

	if !ok {
		return ErrInavlidCredentials
	}

	return nil
}
