package service_test

import (
	"errors"
	"testing"

	"github.com/a5415091-collab/go-gin-todo-app/model"
	"github.com/a5415091-collab/go-gin-todo-app/service"
)

// --- Mock Repository 定義 ---
type MockTodoRepository struct {
	FindAllFunc  func(userID uint) ([]model.Todo, error)
	FindByIDFunc func(userID uint, id uint) (*model.Todo, error)
	CreateFunc   func(todo *model.Todo) (*model.Todo, error)
	UpdateFunc   func(todo *model.Todo) (*model.Todo, error)
	DeleteFunc   func(userID uint, id uint) error
}

func (m *MockTodoRepository) FindAll(userID uint) ([]model.Todo, error) {
	return m.FindAllFunc(userID)
}

func (m *MockTodoRepository) FindByID(userID uint, id uint) (*model.Todo, error) {
	return m.FindByIDFunc(userID, id)
}

func (m *MockTodoRepository) Create(todo *model.Todo) (*model.Todo, error) {
	return m.CreateFunc(todo)
}

func (m *MockTodoRepository) Update(todo *model.Todo) (*model.Todo, error) {
	return m.UpdateFunc(todo)
}

func (m *MockTodoRepository) Delete(userID uint, id uint) error {
	return m.DeleteFunc(userID, id)
}

// --- FindAll ---
func TestTodoService_FindAll(t *testing.T) {

	tests := []struct {
		name       string
		userID     uint
		mockFind   func(userID uint) ([]model.Todo, error)
		expectErr  bool
		expectSize int
	}{
		{
			name:   "success find all",
			userID: 1,
			mockFind: func(userID uint) ([]model.Todo, error) {
				t1 := model.Todo{UserID: userID, Title: "task1", Done: false}
				t1.ID = 1
				t2 := model.Todo{UserID: userID, Title: "task2", Done: true}
				t2.ID = 2
				return []model.Todo{t1, t2}, nil
			},
			expectErr:  false,
			expectSize: 2,
		},
		{
			name:   "failed find all",
			userID: 1,
			mockFind: func(userID uint) ([]model.Todo, error) {
				return nil, errors.New("db error")
			},
			expectErr:  true,
			expectSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockTodoRepository{
				FindAllFunc: tt.mockFind,
			}

			svc := service.NewTodoService(mockRepo)

			result, err := svc.FindAll(tt.userID)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(result) != tt.expectSize {
				t.Errorf("expected %d items, got %d", tt.expectSize, len(result))
			}
		})
	}
}

// --- FindByID ---
func TestTodoService_FindByID(t *testing.T) {

	tests := []struct {
		name         string
		userID       uint
		id           uint
		mockFind     func(userID uint, id uint) (*model.Todo, error)
		expectErr    bool
		expectuserID uint
		expectTitle  string
	}{
		{
			name:   "success find by id",
			userID: 1,
			id:     2,
			mockFind: func(userID uint, id uint) (*model.Todo, error) {
				todo := &model.Todo{UserID: userID, Title: "task2", Done: true}
				todo.ID = id
				return todo, nil
			},
			expectErr:    false,
			expectuserID: 1,
			expectTitle:  "task2",
		},
		{
			name:   "failed findByID",
			userID: 1,
			id:     10,
			mockFind: func(userID uint, id uint) (*model.Todo, error) {
				return nil, errors.New("db error")
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockTodoRepository{
				FindByIDFunc: tt.mockFind,
			}

			svc := service.NewTodoService(mockRepo)

			result, err := svc.FindByID(tt.userID, tt.id)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// 成功時だけ検証
			if !tt.expectErr {
				if result.UserID != tt.expectuserID {
					t.Errorf("userID mismatch: expected %d, got %d", tt.expectuserID, result.UserID)
				}
				if result.Title != tt.expectTitle {
					t.Errorf("title mismatch: expected %s, got %s", tt.expectTitle, result.Title)
				}
			}
		})
	}
}

func TestTodoService_Create(t *testing.T) {

	tests := []struct {
		name     string
		title    string
		mockFunc func(todo *model.Todo) (*model.Todo, error)
		wantErr  bool
	}{
		{
			name:  "success",
			title: "Test Todo",
			mockFunc: func(todo *model.Todo) (*model.Todo, error) {
				todo.ID = 1
				return todo, nil
			},
			wantErr: false,
		},
		{
			name:  "db error",
			title: "Test Todo",
			mockFunc: func(todo *model.Todo) (*model.Todo, error) {
				return nil, errors.New("DB error")
			},
			wantErr: true,
		},
		{
			name:     "empty title should error",
			title:    "",
			mockFunc: nil, // 呼ばれない
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockTodoRepository{
				CreateFunc: func(todo *model.Todo) (*model.Todo, error) {
					if tt.mockFunc != nil {
						return tt.mockFunc(todo)
					}
					return nil, nil
				},
			}

			svc := service.NewTodoService(mockRepo)

			_, err := svc.Create(1, tt.title)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestTodoService_Update(t *testing.T) {

	tests := []struct {
		name       string
		userID     uint
		id         uint
		title      string
		done       *bool
		mockFind   func(userID uint, id uint) (*model.Todo, error)
		mockUpdate func(todo *model.Todo) (*model.Todo, error)
		expectErr  bool
	}{
		{
			name:   "success update",
			userID: 1,
			id:     1,
			title:  "updated title",
			done:   ptrBool(true),
			mockFind: func(userID uint, id uint) (*model.Todo, error) {
				return &model.Todo{
					UserID: userID,
					Title:  "old",
					Done:   false,
				}, nil
			},
			mockUpdate: func(todo *model.Todo) (*model.Todo, error) {
				return todo, nil
			},
			expectErr: false,
		},
		{
			name:   "todo not found",
			userID: 1,
			id:     999,
			title:  "something",
			done:   ptrBool(false),
			mockFind: func(userID uint, id uint) (*model.Todo, error) {
				return nil, errors.New("not found")
			},
			mockUpdate: func(todo *model.Todo) (*model.Todo, error) {
				return nil, nil
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockTodoRepository{
				FindByIDFunc: tt.mockFind,
				UpdateFunc:   tt.mockUpdate,
			}

			svc := service.NewTodoService(mockRepo)

			result, err := svc.Update(tt.userID, tt.id, tt.title, tt.done)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
			if !tt.expectErr && result.Title != tt.title {
				t.Errorf("title mismatch: expected %s, got %s", tt.title, result.Title)
			}
		})
	}
}

func ptrBool(b bool) *bool {
	return &b
}

// --- Delete ---
func TestTodoService_Delete(t *testing.T) {

	tests := []struct {
		name       string
		userID     uint
		id         uint
		mockDelete func(userID uint, id uint) error
		expectErr  bool
	}{
		{
			name:   "success delete",
			userID: 1,
			id:     10,
			mockDelete: func(userID uint, id uint) error {
				return nil
			},
			expectErr: false,
		},
		{
			name:   "failed delete",
			userID: 1,
			id:     10,
			mockDelete: func(userID uint, id uint) error {
				return errors.New("delete failed")
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockTodoRepository{
				DeleteFunc: tt.mockDelete,
			}

			svc := service.NewTodoService(mockRepo)

			err := svc.Delete(tt.userID, tt.id)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
