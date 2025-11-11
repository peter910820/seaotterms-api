package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func selfBrandRouter(blogGroup fiber.Router, dbs map[string]*gorm.DB, dbName string, store *session.Store) {
	selfBrandGroup := blogGroup.Group("/galgame-brands")

	selfBrandGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryAllGalgameBrand(c, dbs[dbName])
	})
	selfBrandGroup.Get("/:brand", func(c *fiber.Ctx) error {
		return api.QueryGalgameBrand(c, dbs[dbName])
	})
	selfBrandGroup.Post("/", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.CreateGalgameBrand(c, dbs[dbName])
	})
	selfBrandGroup.Patch("/:brand", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.UpdateGalgameBrand(c, dbs[dbName])
	})
}
