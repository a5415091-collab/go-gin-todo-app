package repository

import (
	"github.com/a5415091-collab/go-gin-todo-app/db"
	"github.com/a5415091-collab/go-gin-todo-app/model"
)

type TodoRepository interface {
	FindAll() ([]model.Todo, error)
	FindByID(id uint) (*model.Todo, error)
	Create(todo *model.Todo) error
	Update(todo *model.Todo) error
	Delete(id uint) error
}

type todoRepository struct{}

func NewTodoRepository() TodoRepository {
	return &todoRepository{}
}

func (r *todoRepository) FindAll() ([]model.Todo, error) {
	var todos []model.Todo
	result := db.DB.Find(&todos)
	return todos, result.Error
}

func (r *todoRepository) FindByID(id uint) (*model.Todo, error) {
	var todo model.Todo
	result := db.DB.First(&todo, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &todo, nil
}

func (r *todoRepository) Create(todo *model.Todo) error {
	result := db.DB.Create(todo)
	return result.Error
}

func (r *todoRepository) Update(todo *model.Todo) error {
	result := db.DB.Save(todo)
	return result.Error
}

func (r *todoRepository) Delete(id uint) error {
	result := db.DB.Delete(&model.Todo{}, id)
	return result.Error
}
