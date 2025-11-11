package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	api "seaotterms-api/api/blog"
	middleware "seaotterms-api/middleware/blog"
)

// this router is use to check identity for front-end routes
func authRouter(blogGroup fiber.Router, dbs map[string]*gorm.DB, dbName string, store *session.Store) {
	authGroup := blogGroup.Group("/auth")

	// get user info
	authGroup.Get("/", middleware.GetUserInfo(store), func(c *fiber.Ctx) error {
		return api.Auth(c, store)
	})

	authGroup.Post("/login", func(c *fiber.Ctx) error {
		return api.Login(c, store, dbs[dbName])
	})
}
