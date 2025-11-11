package blog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	model "seaotterms-api/model/galgame"
	utils "seaotterms-api/utils/blog"
)

func QueryGameRecord(c *fiber.Ctx, db *gorm.DB) error {
	var responseData []model.PlayRecord

	err := db.Order("COALESCE(updated_at, created_at) DESC").Find(&responseData).Error
	if err != nil {
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	logrus.Info("個人攻略Galgame攻略資料查詢成功")
	response := utils.ResponseFactory(c, fiber.StatusOK, "個人攻略Galgame攻略資料查詢成功", &responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
