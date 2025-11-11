package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	api "seaotterms-api/api/blog"
)

func gameRouter(blogGroup fiber.Router, dbs map[string]*gorm.DB, dbName string, store *session.Store) {
	brandGroup := blogGroup.Group("/galgames")

	brandGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryGame(c, dbs[dbName])
	})

	brandGroup.Post("/", func(c *fiber.Ctx) error {
		return api.CreateGame(c, dbs[dbName])
	})
}
