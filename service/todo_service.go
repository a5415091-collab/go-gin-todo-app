package service

import (
	"errors"

	"github.com/a5415091-collab/go-gin-todo-app/model"
	"github.com/a5415091-collab/go-gin-todo-app/repository"
)

type TodoService interface {
	FindAll(userID uint) ([]model.Todo, error)
	FindByID(userID uint, id uint) (*model.Todo, error)
	Create(userID uint, title string) (*model.Todo, error)
	Update(userID uint, id uint, title string, done bool) (*model.Todo, error)
	Delete(userID uint, id uint) error
}

type todoService struct {
	todoRepo repository.TodoRepository
}

func NewTodoService(todoRepo repository.TodoRepository) TodoService {
	return &todoService{todoRepo}
}

// --- FindAll ---
func (s *todoService) FindAll(userID uint) ([]model.Todo, error) {
	return s.todoRepo.FindAll(userID)
}

// --- FindByID ---
func (s *todoService) FindByID(userID uint, id uint) (*model.Todo, error) {
	todo, err := s.todoRepo.FindByID(userID, id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

// --- Create ---
func (s *todoService) Create(userID uint, title string) (*model.Todo, error) {
	todo := &model.Todo{
		Title:  title,
		UserID: userID,
		Done:   false,
	}
	return s.todoRepo.Create(todo)
}

// --- Update ---
func (s *todoService) Update(userID uint, id uint, title string, done bool) (*model.Todo, error) {
	todo, err := s.todoRepo.FindByID(userID, id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	todo.Title = title
	todo.Done = done
	return s.todoRepo.Update(todo)
}

// --- Delete ---
func (s *todoService) Delete(userID uint, id uint) error {
	return s.todoRepo.Delete(userID, id)
}
