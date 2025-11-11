package blog

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	dto "seaotterms-api/dto/blog"
	model "seaotterms-api/model/blog"
	utils "seaotterms-api/utils/blog"
)

func QueryTodoByOwner(c *fiber.Ctx, db *gorm.DB) error {
	// URL decoding
	owner, err := url.QueryUnescape(c.Params("owner"))
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, "客戶端資料錯誤", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var responseData []model.Todo
	err = db.Where("owner = ?", owner).Order("created_at DESC").Find(&responseData).Error
	if err != nil {
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			logrus.Error(err)
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, err.Error(), nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	logrus.Infof("查詢%s的Todo資料成功", owner)
	response := utils.ResponseFactory(c, fiber.StatusOK, fmt.Sprintf("查詢%s的Todo資料成功", owner), &responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateTodo(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData model.Todo
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, "客戶端資料錯誤", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// handle topic
	lastSlashIndex := strings.LastIndex(clientData.Topic, "/")
	if lastSlashIndex != -1 {
		clientData.Topic = clientData.Topic[:lastSlashIndex]
	} else {
		logrus.Error("topic value has error")
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, "客戶端資料轉換錯誤", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	data := model.Todo{
		ID:         clientData.ID,
		Owner:      clientData.Owner,
		Topic:      clientData.Topic,
		Title:      clientData.Title,
		Status:     clientData.Status,
		Deadline:   clientData.Deadline,
		CreateName: clientData.CreateName,
		UpdateName: clientData.UpdateName,
	}

	err := db.Create(&data).Error
	if err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	responseData, err := getTodo(c, db)
	if err != nil {
		logrus.Error(err)
		if err == gorm.ErrRecordNotFound {
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, "找不到該Todo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	logrus.Infof("資料 %s 創建成功", clientData.Title)
	response := utils.ResponseFactory(c, fiber.StatusOK, fmt.Sprintf("資料 %s 創建成功", clientData.Title), responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func UpdateTodoStatus(c *fiber.Ctx, db *gorm.DB) error {

	// load client data
	var clientData dto.TodoUpdateRequest
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, "客戶端資料錯誤", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	clientData.UpdatedAt = time.Now()
	err := db.Model(&model.Todo{}).Where("id = ?", c.Params("id")).
		Select("status", "updated_at", "update_name").
		Updates(clientData).Error
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, "找不到該Todo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	responseData, err := getTodo(c, db)
	if err != nil {
		logrus.Error(err)
		if err == gorm.ErrRecordNotFound {
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, "找不到該Todo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	logrus.Infof("Todo %s 更新成功", c.Params("id"))
	response := utils.ResponseFactory(c, fiber.StatusOK, fmt.Sprintf("Todo %s 更新成功", c.Params("id")), responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteTodo(c *fiber.Ctx, db *gorm.DB) error {
	err := db.Where("id = ?", c.Params("id")).Delete(&model.Todo{}).Error
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, "找不到該Todo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	responseData, err := getTodo(c, db)
	if err != nil {
		logrus.Error(err)
		if err == gorm.ErrRecordNotFound {
			response := utils.ResponseFactory[any](c, fiber.StatusNotFound, "找不到該Todo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	logrus.Infof("Todo %s 刪除成功", c.Params("id"))
	response := utils.ResponseFactory(c, fiber.StatusOK, fmt.Sprintf("Todo %s 刪除成功", c.Params("id")), responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

// 用使用者登入資料取得該使用者的全部Todo資料
// 用在增、改、刪三個API的回傳值，降低前端Request的次數
func getTodo(c *fiber.Ctx, db *gorm.DB) (*[]model.Todo, error) {
	userInfo, ok := c.Locals("user_info").(*dto.UserInfo)
	if !ok {
		logrus.Fatal("使用者登入版號表異常")
	}
	var responseData []model.Todo
	err := db.Where("owner = ?", userInfo.Username).Order("created_at DESC").Find(&responseData).Error
	if err != nil {
		return nil, err
	}
	return &responseData, nil
}
