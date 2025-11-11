package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func articleRouter(blogGroup fiber.Router, dbs map[string]*gorm.DB, dbName string, store *session.Store) {
	articleGroup := blogGroup.Group("/articles")

	articleGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryArticle(c, dbs[dbName])
	})

	articleGroup.Get("/:id", func(c *fiber.Ctx) error {
		return api.QueryArticle(c, dbs[dbName])
	})

	// No middleware has been implemented yet
	articleGroup.Post("/", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.CreateArticle(c, dbs[dbName])
	})

	// articleGroup.Post("/:id", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
	// 	return api.ModifyArticle(c, dbs[dbName])
	// })

	articleGroup.Delete("/:id", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.DeleteArticle(c, dbs[dbName])
	})
}
