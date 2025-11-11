package model

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"

	blogmodel "seaotterms-api/model/blog"
	galmodel "seaotterms-api/model/galgame"
	teachmodel "seaotterms-api/model/teach"
)

func Migration(dbName string, db *gorm.DB) {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file error: %v", err)
	}

	switch dbName {
	case os.Getenv("DATABASE_NAME"):
		db.AutoMigrate(&teachmodel.Series{})
		db.AutoMigrate(&teachmodel.Article{})
		db.AutoMigrate(&teachmodel.Comment{})
	case os.Getenv("DATABASE_NAME2"):
		// db.AutoMigrate(&galmodel.DownloadArticle{})
		db.AutoMigrate(&galmodel.Brand{})
		db.AutoMigrate(&galmodel.Game{})
		db.AutoMigrate(&galmodel.PlayRecord{})
		db.AutoMigrate(&galmodel.SelfBrand{}) // old
		db.AutoMigrate(&galmodel.SelfGame{})  //old
		db.AutoMigrate(&galmodel.User{})
		db.AutoMigrate(&galmodel.Tag{})
		db.AutoMigrate(&galmodel.Article{})
		db.AutoMigrate(&galmodel.Log{})
		db.AutoMigrate(&galmodel.TmpData{})
	case os.Getenv("DATABASE_NAME3"):
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
