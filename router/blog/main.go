package blog

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	seaottermsdb "seaotterms-db"

	middleware "seaotterms-api/middleware/blog"
)

// 除了身份驗證表的資料庫，其餘資料庫名稱都定義在各站台router包的main.go中
func BlogRouter(apiGroup fiber.Router, dbm *seaottermsdb.DBModel, store *session.Store) {
	blogGroup := apiGroup.Group("/blog")

	dbName := os.Getenv("DATABASE_NAME2")

	blogGroup.Use(middleware.GetUserInfo(store)) // global middleware

	// article
	articleRouter(blogGroup, dbm, dbName, store)
	tagRouter(blogGroup, dbm, dbName, store)

	todoRouter(blogGroup, dbm, dbName, store)
	systemTodoRouter(blogGroup, dbm, dbName, store)
	userRouter(blogGroup, dbm, dbName, store)
	todoTopicRouter(blogGroup, dbm, dbName, store)

	// auth
	authRouter(blogGroup, dbm, dbName, store)
}
