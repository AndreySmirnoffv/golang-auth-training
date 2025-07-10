package usecases

import (
	"errors"

	"github.com/AndreySmirnoffv/golang-auth-training/internal/entity"
)

var ErrEmailExists = errors.New("email already registered")

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

	return uc.repo.Save(u)
}
