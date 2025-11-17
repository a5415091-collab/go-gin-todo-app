package db

import (
	"log"

	"github.com/glebarez/sqlite" // ← これが modernc ベースのドライバ
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error

	DB, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
}
