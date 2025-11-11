package blog

import (
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	dto "seaotterms-api/dto/blog"
	model "seaotterms-api/model/blog"
	utils "seaotterms-api/utils/blog"
)

func QueryTodoTopic(c *fiber.Ctx, db *gorm.DB) error {
	// URL decoding
	owner, err := url.QueryUnescape(c.Params("owner"))
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, "客戶端資料錯誤", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	data := []model.TodoTopic{}
	err = db.Where("topic_owner = ?", owner).Order("topic_name DESC").Find(&data).Error
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, "找不到TodoTopic資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	logrus.Info("查詢TodoTopic資料成功")
	response := utils.ResponseFactory(c, fiber.StatusOK, "查詢TodoTopic資料成功", &data)
	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateTodoTopic(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData dto.TodoTopicCreateRequest
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	data := model.TodoTopic{
		TopicName:  clientData.TopicName,
		TopicOwner: clientData.TopicOwner,
		UpdateName: clientData.UpdateName,
	}

	err := db.Create(&data).Error
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	logrus.Infof("資料 %s 創建成功", clientData.TopicName)
	response := utils.ResponseFactory[any](c, fiber.StatusOK, fmt.Sprintf("資料 %s 創建成功", clientData.TopicName), nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
