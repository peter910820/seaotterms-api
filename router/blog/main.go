package blog

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	middleware "seaotterms-api/middleware/blog"
)

// 除了身份驗證表的資料庫，其餘資料庫名稱都定義在各站台router包的main.go中
func BlogRouter(apiGroup fiber.Router, dbs map[string]*gorm.DB, store *session.Store) {
	blogGroup := apiGroup.Group("/blog")

	dbName := os.Getenv("DATABASE_NAME2")

	blogGroup.Use(middleware.GetUserInfo(store)) // global middleware

	// article
	articleRouter(blogGroup, dbs, dbName, store)
	tagRouter(blogGroup, dbs, dbName, store)

	todoRouter(blogGroup, dbs, dbName, store)
	systemTodoRouter(blogGroup, dbs, dbName, store)
	userRouter(blogGroup, dbs, dbName, store)
	todoTopicRouter(blogGroup, dbs, dbName, store)

	// auth
	authRouter(blogGroup, dbs, dbName, store)
}
