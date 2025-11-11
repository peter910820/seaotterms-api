package blog

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	dto "seaotterms-api/dto/blog"
	model "seaotterms-api/model/blog"
	utils "seaotterms-api/utils/blog"
)

func QueryTag(c *fiber.Ctx, db *gorm.DB) error {
	var responseData []model.Tag

	err := db.Order("created_at desc").Find(&responseData).Error
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, err.Error(), nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	logrus.Info("Tag資料查詢成功")
	response := utils.ResponseFactory(c, fiber.StatusOK, "Tag資料查詢成功", &responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateTag(c *fiber.Ctx, db *gorm.DB) error {
	var clientData dto.TagCreateRequest

	// load client data
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	dataCreate := model.Tag{
		Name: clientData.Name,
	}
	if strings.TrimSpace(clientData.IconName) != "" {
		dataCreate.IconName = clientData.IconName
	}
	err := db.Create(&dataCreate).Error
	if err != nil {
		logrus.Println("錯誤:", err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	response := utils.ResponseFactory[any](c, fiber.StatusOK, "成功建立Tag", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
