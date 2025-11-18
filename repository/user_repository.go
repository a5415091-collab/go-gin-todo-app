package repository

import (
	"github.com/a5415091-collab/go-gin-todo-app/db"
	"github.com/a5415091-collab/go-gin-todo-app/model"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	result := db.DB.Create(user)
	return result.Error
}
