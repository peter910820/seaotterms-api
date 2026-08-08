package blog

import (
	seaottermsdb "seaotterms-db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func todoTopicRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, dbName string, store *session.Store) {
	todoTopicGroup := blogGroup.Group("/todo-topics")

	todoTopicGroup.Get("/:owner", func(c *fiber.Ctx) error {
		return api.QueryTodoTopic(c, dbm.DB)
	})
	todoTopicGroup.Post("/", middleware.CheckLogin(store), func(c *fiber.Ctx) error {
		return api.CreateTodoTopic(c, dbm.DB)
	})
}
