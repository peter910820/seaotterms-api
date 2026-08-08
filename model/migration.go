package model

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"

	blogmodel "seaotterms-api/model/blog"
)

func Migration(dbName string, db *gorm.DB) {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file error: %v", err)
	}

	switch dbName {
	case os.Getenv("DATABASE_NAME"):
		db.AutoMigrate(&blogmodel.User{})
		db.AutoMigrate(&blogmodel.Tag{})
		db.AutoMigrate(&blogmodel.Article{})
		db.AutoMigrate(&blogmodel.Todo{})
		db.AutoMigrate(&blogmodel.TodoTopic{})
		db.AutoMigrate(&blogmodel.SystemTodo{})
	default:
		logrus.Fatal("error in migration function")
	}
}
