package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func tagRouter(blogGroup fiber.Router, dbs map[string]*gorm.DB, dbName string, store *session.Store) {
	tagGroup := blogGroup.Group("/tags")

	tagGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryTag(c, dbs[dbName])
	})

	tagGroup.Get("/:name", func(c *fiber.Ctx) error {
		return api.QueryArticleForTag(c, dbs[dbName])
	})

	// create tag
	tagGroup.Post("/", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.CreateTag(c, dbs[dbName])
	})
}
