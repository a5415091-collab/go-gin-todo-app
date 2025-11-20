package service

import (
	"errors"

	"github.com/a5415091-collab/go-gin-todo-app/model"
	"github.com/a5415091-collab/go-gin-todo-app/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(email, password string) error
	Login(email, password string) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

// Signup
func (s *authService) Signup(email, password string) error {
	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:    email,
		Password: string(hashed),
	}

	return s.userRepo.Create(user)

}

// Login
func (s *authService) Login(email, password string) (*model.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
