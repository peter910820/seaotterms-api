package blog

import (
	seaottermsdb "seaotterms-db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

func userRouter(blogGroup fiber.Router, dbm *seaottermsdb.DBModel, dbName string, store *session.Store) {
	userGroup := blogGroup.Group("/users")

	userGroup.Get("/", middleware.CheckManagement(store), func(c *fiber.Ctx) error {
		return api.QueryUser(c, dbm.DB)
	})

	userGroup.Post("/", func(c *fiber.Ctx) error {
		return api.CreateUser(c, dbm.DB)
	})
	userGroup.Patch("/:id", middleware.CheckLogin(store), func(c *fiber.Ctx) error {
		return api.UpdateUser(c, dbm.DB, store)
	})
}
