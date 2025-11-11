package blog

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	dto "seaotterms-api/dto/blog"
	model "seaotterms-api/model/galgame"
	utils "seaotterms-api/utils/blog"
)

func QueryGame(c *fiber.Ctx, db *gorm.DB) error {
	var responseData []model.Game

	err := db.Order("COALESCE(updated_at, created_at) DESC").Find(&responseData).Error
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	logrus.Info("Game資料查詢成功")
	response := utils.ResponseFactory(c, fiber.StatusOK, "Game資料查詢成功", &responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateGame(c *fiber.Ctx, db *gorm.DB) error {
	var requestData dto.GameCreateRequest

	if err := c.BodyParser(&requestData); err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err := db.First(&model.Game{}, requestData.BrandID).Error
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, "找不到Brand資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	data := model.Game{
		Name:            requestData.Name,
		ChineseName:     requestData.ChineseName,
		BrandID:         requestData.BrandID,
		AllAges:         requestData.AllAges,
		ReleaseDate:     requestData.ReleaseDate,
		OpUrl:           requestData.OpUrl,
		GameDescription: requestData.GameDescription,
		CreatedAt:       time.Now(),
		CreatedName:     "seaotterms",
	}

	if err := db.Create(&data).Error; err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	logrus.Info("Galgame資料建立成功: " + requestData.Name)
	response := utils.ResponseFactory[any](c, fiber.StatusOK, "Galgame資料建立成功: "+requestData.Name, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
