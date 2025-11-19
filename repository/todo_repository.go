package repository

import (
	"github.com/a5415091-collab/go-gin-todo-app/db"
	"github.com/a5415091-collab/go-gin-todo-app/model"
)

type TodoRepository interface {
	FindAll(userID uint) ([]model.Todo, error)
	FindByID(userID uint, id uint) (*model.Todo, error)
	Create(todo *model.Todo) (*model.Todo, error)
	Update(todo *model.Todo) (*model.Todo, error)
	Delete(userID uint, id uint) error
}

type todoRepository struct{}

func NewTodoRepository() TodoRepository {
	return &todoRepository{}
}

func (r *todoRepository) FindAll(userID uint) ([]model.Todo, error) {
	var todos []model.Todo
	err := db.DB.Where("user_id = ?", userID).Find(&todos).Error
	return todos, err
}

func (r *todoRepository) FindByID(userID uint, id uint) (*model.Todo, error) {
	var todo model.Todo
	err := db.DB.Where("user_id = ? AND id = ?", userID, id).First(&todo).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) Create(todo *model.Todo) (*model.Todo, error) {
	result := db.DB.Create(todo)
	if result.Error != nil {
		return nil, result.Error
	}
	return todo, nil
}

func (r *todoRepository) Update(todo *model.Todo) (*model.Todo, error) {
	result := db.DB.
		Where("id = ? AND user_id = ?", todo.ID, todo.UserID).
		Updates(todo)

	if result.Error != nil {
		return nil, result.Error
	}
	return todo, nil
}

func (r *todoRepository) Delete(userID uint, id uint) error {
	result := db.DB.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Todo{})

	return result.Error
}
