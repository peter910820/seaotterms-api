package gal

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	api "seaotterms-api/api/gal"
)

func GalRouter(apiGroup fiber.Router, dbs map[string]*gorm.DB) {
	galGroup := apiGroup.Group("/gal")
	dbName := os.Getenv("DATABASE_NAME2")

	authRouter(galGroup, dbs, dbName)
	loginRouter(galGroup, dbs, dbName)
}

func authRouter(apiGroup fiber.Router, dbs map[string]*gorm.DB, dbName string) {
	apiGroup.Post("/register", func(c *fiber.Ctx) error {
		return api.Register(c, dbs[dbName])
	})

	apiGroup.Get("/register/:mail_name/:register_key", func(c *fiber.Ctx) error {
		return api.RegisterKeyCheck(c, dbs[dbName])
	})
}

func loginRouter(apiGroup fiber.Router, dbs map[string]*gorm.DB, dbName string) {
	apiGroup.Post("/login", func(c *fiber.Ctx) error {
		return api.Register(c, dbs[dbName])
	})
}
