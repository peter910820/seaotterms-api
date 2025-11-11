package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func todoRouter(blogGroup fiber.Router, dbs map[string]*gorm.DB, dbName string, store *session.Store) {
	todoGroup := blogGroup.Group("/todos")

	todoGroup.Get("/:owner", func(c *fiber.Ctx) error {
		return api.QueryTodoByOwner(c, dbs[dbName])
	})
	todoGroup.Post("/", middleware.CheckLogin(store), func(c *fiber.Ctx) error {
		return api.CreateTodo(c, dbs[dbName])
	})
	todoGroup.Patch("/:id", middleware.CheckLogin(store), func(c *fiber.Ctx) error {
		return api.UpdateTodoStatus(c, dbs[dbName])
	})
	todoGroup.Delete("/:id", middleware.CheckLogin(store), func(c *fiber.Ctx) error {
		return api.DeleteTodo(c, dbs[dbName])
	})
}
