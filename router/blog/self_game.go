package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func selfGameRouter(blogGroup fiber.Router, dbs map[string]*gorm.DB, dbName string, store *session.Store) {
	selfGameGroup := blogGroup.Group("/galgames")

	selfGameGroup.Get("/s/:name", func(c *fiber.Ctx) error {
		return api.QueryGalgame(c, dbs[dbName])
	})
	selfGameGroup.Get("/:brand", func(c *fiber.Ctx) error {
		return api.QueryGalgameByBrand(c, dbs[dbName])
	})
	selfGameGroup.Patch("/develop/:name", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.UpdateGalgameDevelop(c, dbs[dbName])
	})
	selfGameGroup.Post("/", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.CreateGalgame(c, dbs[dbName])
	})
}
