package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	api "seaotterms-api/api/blog"
)

func brandRouter(blogGroup fiber.Router, dbs map[string]*gorm.DB, dbName string, store *session.Store) {
	brandGroup := blogGroup.Group("/brands")

	brandGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryBrand(c, dbs[dbName])
	})

	brandGroup.Post("/", func(c *fiber.Ctx) error {
		return api.CreateBrand(c, dbs[dbName])
	})
}
