package teach

import (
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms-db/teach"
)

// Request
type (
	SeriesCreateRequest struct {
		Title string `gorm:"NOT NULL" json:"title"`
		Image string `json:"image"`
	}

	SeriesModifyRequest struct {
		Title string `gorm:"NOT NULL" json:"title"`
		Image string `json:"image"`
	}
)

func querySeries(c *fiber.Ctx, db *gorm.DB) error {
	var responseData []teach.Series
	// URL decoding
	id, err := url.QueryUnescape(c.Params("id"))
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if id == "" {
		responseData, err = teach.FindAllSeries(db)
		if err != nil {
			return dbErrHandler(c, err)
		}
	} else {
		data, err := teach.FindSeriesByID(db, id)
		if err != nil {
			return dbErrHandler(c, err)
		}
		responseData = append(responseData, *data)

	}

	logrus.Info("Query series table success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "查詢Series資料成功",
		"data": responseData,
	})
}

func createSeries(c *fiber.Ctx, db *gorm.DB) error {
	var clientData SeriesCreateRequest
	// load client data
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	data := teach.Series{
		Title:       clientData.Title,
		Image:       clientData.Image,
		CreatedAt:   time.Now(),
		CreatedName: "Root",
	}

	err := teach.CreateSeries(db, &data)
	if err != nil {
		return dbErrHandler(c, err)
	}

	logrus.Print("新增Series資料成功")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "新增Series資料成功",
		"data": data,
	})
}

func modifySeries(c *fiber.Ctx, db *gorm.DB) error {
	var clientData SeriesModifyRequest
	// load client data
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	data := teach.Series{
		Title:       clientData.Title,
		Image:       clientData.Image,
		UpdatedAt:   time.Now(),
		UpdatedName: "Root",
	}

	err := teach.UpdateSeries(db, c.Params("id"), &data)
	if err != nil {
		return dbErrHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "修改Series資料成功",
		"data": data,
	})
}
