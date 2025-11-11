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

func QueryBrand(c *fiber.Ctx, db *gorm.DB) error {
	var responseData []model.Brand

	err := db.Order("COALESCE(updated_at, created_at) DESC").Find(&responseData).Error
	if err != nil {
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	logrus.Info("Brand資料查詢成功")
	response := utils.ResponseFactory(c, fiber.StatusOK, "Brand資料查詢成功", &responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateBrand(c *fiber.Ctx, db *gorm.DB) error {
	var requestData dto.BrandCreateRequest

	if err := c.BodyParser(&requestData); err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	data := model.Brand{
		Name:        requestData.Name,
		WorkAmount:  requestData.WorkAmount,
		OfficialUrl: requestData.OfficialUrl,
		Dissolution: requestData.Dissolution,
		CreatedAt:   time.Now(),
		CreatedName: "seaotterms",
	}

	if err := db.Create(&data).Error; err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	logrus.Info("Galgame品牌資料建立成功: " + requestData.Name)
	response := utils.ResponseFactory[any](c, fiber.StatusOK, "Galgame品牌資料建立成功: "+requestData.Name, nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
