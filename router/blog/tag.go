package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	seaottermsdb "seaotterms-db"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func tagRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, dbName string, store *session.Store) {
	tagGroup := blogGroup.Group("/tags")

	tagGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryTag(c, dbm.DB)
	})

	tagGroup.Get("/:name", func(c *fiber.Ctx) error {
		return api.QueryArticleForTag(c, dbm.DB)
	})

	// create tag
	tagGroup.Post("/", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.CreateTag(c, dbm.DB)
	})
}
