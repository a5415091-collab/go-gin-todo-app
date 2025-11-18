package service

import (
	"errors"

	"github.com/a5415091-collab/go-gin-todo-app/model"
	"github.com/a5415091-collab/go-gin-todo-app/repository"
)

type TodoService interface {
	FindAll() ([]model.Todo, error)
	FindByID(id uint) (*model.Todo, error)
	Create(title string, userID uint) error
	Update(id uint, title string, done bool) error
	Delete(id uint) error
}

type todoService struct {
	todoRepo repository.TodoRepository
}

func NewTodoService(todoRepo repository.TodoRepository) TodoService {
	return &todoService{todoRepo}
}

// --- 既存の FindAll ---
func (s *todoService) FindAll() ([]model.Todo, error) {
	return s.todoRepo.FindAll()
}

// --- FindByID ---
func (s *todoService) FindByID(id uint) (*model.Todo, error) {
	todo, err := s.todoRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

// --- Create ---
func (s *todoService) Create(title string, userID uint) error {
	todo := &model.Todo{
		Title:  title,
		UserID: userID,
		Done:   false,
	}
	return s.todoRepo.Create(todo)
}

// --- Update ---
func (s *todoService) Update(id uint, title string, done bool) error {
	todo, err := s.todoRepo.FindByID(id)
	if err != nil {
		return errors.New("todo not found")
	}
	todo.Title = title
	todo.Done = done
	return s.todoRepo.Update(todo)
}

// --- Delete ---
func (s *todoService) Delete(id uint) error {
	return s.todoRepo.Delete(id)
}
