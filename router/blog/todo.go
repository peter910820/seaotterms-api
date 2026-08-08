package blog

import (
	seaottermsdb "seaotterms-db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func todoRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, dbName string, store *session.Store) {
	todoGroup := blogGroup.Group("/todos")

	todoGroup.Get("/:owner", func(c *fiber.Ctx) error {
		return api.QueryTodoByOwner(c, dbm.DB)
	})
	todoGroup.Post("/", middleware.CheckLogin(store), func(c *fiber.Ctx) error {
		return api.CreateTodo(c, dbm.DB)
	})
	todoGroup.Patch("/:id", middleware.CheckLogin(store), func(c *fiber.Ctx) error {
		return api.UpdateTodoStatus(c, dbm.DB)
	})
	todoGroup.Delete("/:id", middleware.CheckLogin(store), func(c *fiber.Ctx) error {
		return api.DeleteTodo(c, dbm.DB)
	})
}
