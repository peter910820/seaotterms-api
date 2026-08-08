package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	seaottermsdb "seaotterms-db"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func articleRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, dbName string, store *session.Store) {
	articleGroup := blogGroup.Group("/articles")

	articleGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryArticle(c, dbm.DB)
	})

	articleGroup.Get("/:id", func(c *fiber.Ctx) error {
		return api.QueryArticle(c, dbm.DB)
	})

	// No middleware has been implemented yet
	articleGroup.Post("/", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.CreateArticle(c, dbm.DB)
	})

	// articleGroup.Post("/:id", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
	// 	return api.ModifyArticle(c, dbm.DB)
	// })

	articleGroup.Delete("/:id", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.DeleteArticle(c, dbm.DB)
	})
}
