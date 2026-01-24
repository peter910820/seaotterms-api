package teach

import (
	seaottermsdb "seaotterms-db"

	"github.com/gofiber/fiber/v2"
)

func TeachRouter(apiGroup fiber.Router, dbm *seaottermsdb.DBModel) {
	teachGroup := apiGroup.Group("/teach")

	seriesRouter(teachGroup, dbm)
	articleRouter(teachGroup, dbm)
	commentApiRouter(teachGroup, dbm)
}

func seriesRouter(apiGroup fiber.Router, dbm *seaottermsdb.DBModel) {
	apiGroup.Get("/series", func(c *fiber.Ctx) error {
		return querySeries(c, dbm.DB)
	})

	apiGroup.Post("/series", func(c *fiber.Ctx) error {
		return createSeries(c, dbm.DB)
	})

	apiGroup.Patch("/series/:id", func(c *fiber.Ctx) error {
		return modifySeries(c, dbm.DB)
	})
}

func articleRouter(apiGroup fiber.Router, dbm *seaottermsdb.DBModel) {
	apiGroup.Get("/article", func(c *fiber.Ctx) error {
		return queryArticle(c, dbm.DB)
	})

	apiGroup.Post("/article", func(c *fiber.Ctx) error {
		return createArticle(c, dbm.DB)
	})
}

func commentApiRouter(apiGroup fiber.Router, dbm *seaottermsdb.DBModel) {
}
