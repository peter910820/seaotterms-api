package teach

import (
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	dto "seaotterms-api/dto/teach"
	model "seaotterms-api/model/teach"
)

// query series
func QuerySeries(c *fiber.Ctx, db *gorm.DB) error {
	var responseData []model.Series
	var r *gorm.DB
	// URL decoding
	id, err := url.QueryUnescape(c.Params("id"))
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	if id == "" {
		r = db.Order("COALESCE(updated_at, created_at) DESC").Find(&responseData)
	} else {
		r = db.Where("id = ?", id).Order("COALESCE(updated_at, created_at) DESC").Find(&responseData)
	}
	if r.Error != nil {
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			logrus.Error(r.Error)
			//404
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			logrus.Fatal(r.Error.Error())
			// 500
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		}
	}
	logrus.Info("Query series table success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "查詢Series資料成功",
		"data": responseData,
	})
}

func CreateSeries(c *fiber.Ctx, db *gorm.DB) error {
	var clientData dto.SeriesCreateRequest
	// load client data
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	data := model.Series{
		Title:       clientData.Title,
		Image:       clientData.Image,
		CreatedAt:   time.Now(),
		CreatedName: "Root",
	}

	r := db.Create(&data)
	if r.Error != nil {
		logrus.Error(r.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": r.Error.Error(),
		})
	}
	logrus.Print("新增Series資料成功")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "新增Series資料成功",
		"data": data,
	})
}

func ModifySeries(c *fiber.Ctx, db *gorm.DB) error {
	var clientData dto.SeriesModifyRequest
	// load client data
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	data := model.Series{
		Title:       clientData.Title,
		Image:       clientData.Image,
		UpdatedAt:   time.Now(),
		UpdatedName: "Root",
	}

	r := db.Model(&model.Series{}).Where("id = ?", c.Params("id")).
		Select("title", "image", "updated_at", "updated_name").
		Updates(data)
	if r.Error != nil {
		logrus.Error(r.Error)
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "修改Series資料成功",
		"data": data,
	})
}
