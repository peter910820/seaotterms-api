package blog

import (
	seaottermsdb "seaotterms-db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func systemTodoRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, dbName string, store *session.Store) {
	systemTodoGroup := blogGroup.Group("/system-todos")

	systemTodoGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QuerySystemTodo(c, dbm.DB)
	})

	systemTodoGroup.Post("/", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.CreateSystemTodo(c, dbm.DB)
	})

	systemTodoGroup.Patch("/:id", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.UpdateSystemTodo(c, dbm.DB)
	})

	// quick update
	systemTodoGroup.Patch("/quick/:id", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.QuickUpdateSystemTodo(c, dbm.DB)
	})

	systemTodoGroup.Delete("/:id", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.DeleteSystemTodo(c, dbm.DB)
	})
}
