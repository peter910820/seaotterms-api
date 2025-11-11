package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func userRouter(blogGroup fiber.Router, dbs map[string]*gorm.DB, dbName string, store *session.Store) {
	userGroup := blogGroup.Group("/users")

	userGroup.Get("/", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.QueryUser(c, dbs[dbName])
	})

	userGroup.Post("/", func(c *fiber.Ctx) error {
		return api.CreateUser(c, dbs[dbName])
	})
	userGroup.Patch("/:id", middleware.CheckLogin(store), func(c *fiber.Ctx) error {
		return api.UpdateUser(c, dbs[dbName], store)
	})
}
