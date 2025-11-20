package service_test

import (
	"errors"
	"testing"

	"github.com/a5415091-collab/go-gin-todo-app/model"
	"github.com/a5415091-collab/go-gin-todo-app/service"
	"golang.org/x/crypto/bcrypt"
)

// --- Mock Repository ---

type MockUserRepository struct {
	FindByEmailFunc func(email string) (*model.User, error)
	CreateFunc      func(user *model.User) error
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	return m.FindByEmailFunc(email)
}

func (m *MockUserRepository) Create(user *model.User) error {
	return m.CreateFunc(user)
}

// =====================
//
//	Signup Test
//
// =====================
func TestAuthService_Signup(t *testing.T) {

	tests := []struct {
		name       string
		email      string
		password   string
		mockFind   func(email string) (*model.User, error)
		mockCreate func(user *model.User) error
		expectErr  bool
	}{
		{
			name:     "success signup",
			email:    "test@example.com",
			password: "pass1234",
			mockFind: func(email string) (*model.User, error) {
				return nil, nil
			},
			mockCreate: func(user *model.User) error {
				return nil
			},
			expectErr: false,
		},
		{
			name:     "email already exists",
			email:    "test@example.com",
			password: "pass1234",
			mockFind: func(email string) (*model.User, error) {
				u := &model.User{Email: email}
				u.ID = 1
				return u, nil
			},
			mockCreate: func(user *model.User) error {
				return nil
			},
			expectErr: true,
		},
		{
			name:     "db create error",
			email:    "test@example.com",
			password: "pass1234",
			mockFind: func(email string) (*model.User, error) {
				return nil, nil
			},
			mockCreate: func(user *model.User) error {
				return errors.New("db error")
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockUserRepository{
				FindByEmailFunc: tt.mockFind,
				CreateFunc:      tt.mockCreate,
			}

			svc := service.NewAuthService(mockRepo)

			err := svc.Signup(tt.email, tt.password)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// =====================
//
//	Login Test
//
// =====================
func TestAuthService_Login(t *testing.T) {

	hashed, _ := bcrypt.GenerateFromPassword([]byte("pass1234"), bcrypt.DefaultCost)

	tests := []struct {
		name      string
		email     string
		password  string
		mockFind  func(email string) (*model.User, error)
		expectErr bool
	}{
		{
			name:     "success login",
			email:    "test@example.com",
			password: "pass1234",
			mockFind: func(email string) (*model.User, error) {
				u := &model.User{Email: email, Password: string(hashed)}
				u.ID = 1
				return u, nil
			},
			expectErr: false,
		},
		{
			name:     "email not found",
			email:    "none@example.com",
			password: "pass1234",
			mockFind: func(email string) (*model.User, error) {
				return nil, nil
			},
			expectErr: true,
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "wrongpass",
			mockFind: func(email string) (*model.User, error) {
				u := &model.User{Email: email, Password: string(hashed)}
				u.ID = 1
				return u, nil
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockUserRepository{
				FindByEmailFunc: tt.mockFind,
			}

			svc := service.NewAuthService(mockRepo)

			_, err := svc.Login(tt.email, tt.password)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
