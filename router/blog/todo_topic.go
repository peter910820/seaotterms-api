package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func todoTopicRouter(blogGroup fiber.Router, dbs map[string]*gorm.DB, dbName string, store *session.Store) {
	todoTopicGroup := blogGroup.Group("/todo-topics")

	todoTopicGroup.Get("/:owner", func(c *fiber.Ctx) error {
		return api.QueryTodoTopic(c, dbs[dbName])
	})
	todoTopicGroup.Post("/", middleware.CheckLogin(store), func(c *fiber.Ctx) error {
		return api.CreateTodoTopic(c, dbs[dbName])
	})
}
