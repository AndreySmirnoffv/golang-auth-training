package db

import (
	"github.com/AndreySmirnoffv/golang-auth-training/internal/entity"
	"github.com/AndreySmirnoffv/golang-auth-training/internal/usecases"
	"gorm.io/gorm"
)

type UserRepoGORM struct {
	db *gorm.DB
}

func NewUserRepoGORM(db *gorm.DB) usecases.UserRepo {
	return &UserRepoGORM{db: db}
}

func (r *UserRepoGORM) Save(u *entity.User) error {
	model := UserModel{Email: u.Email, Password: u.Password}
	return r.db.Create(&model).Error
}

func (r *UserRepoGORM) FindByEmail(email string) (*entity.User, error) {
	var model UserModel

	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}
	return &entity.User{ID: model.ID, Email: model.Email, Password: model.Password}, nil

}
