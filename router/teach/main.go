package teach

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	api "seaotterms-api/api/teach"
)

func TeachRouter(apiGroup fiber.Router, dbs map[string]*gorm.DB) {
	teachGroup := apiGroup.Group("/teach")
	dbName := os.Getenv("DATABASE_NAME")

	seriesRouter(teachGroup, dbs, dbName)
	articleRouter(teachGroup, dbs, dbName)
	commentApiRouter(teachGroup, dbs, dbName)
}

func seriesRouter(apiGroup fiber.Router, dbs map[string]*gorm.DB, dbName string) {
	apiGroup.Get("/series", func(c *fiber.Ctx) error {
		return api.QuerySeries(c, dbs[dbName])
	})

	apiGroup.Post("/series", func(c *fiber.Ctx) error {
		return api.CreateSeries(c, dbs[dbName])
	})

	apiGroup.Patch("/series/:id", func(c *fiber.Ctx) error {
		return api.ModifySeries(c, dbs[dbName])
	})
}

func articleRouter(apiGroup fiber.Router, dbs map[string]*gorm.DB, dbName string) {
	apiGroup.Get("/article", func(c *fiber.Ctx) error {
		return api.QueryArticle(c, dbs[dbName])
	})

	apiGroup.Post("/article", func(c *fiber.Ctx) error {
		return api.CreateArticle(c, dbs[dbName])
	})
}

func commentApiRouter(apiGroup fiber.Router, dbs map[string]*gorm.DB, dbName string) {
}
