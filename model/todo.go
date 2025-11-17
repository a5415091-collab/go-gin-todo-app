package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	UserID uint
	Title  string
	Done   bool
}
