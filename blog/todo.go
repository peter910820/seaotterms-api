package blog

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	model "seaotterms-db/blog"
)

type TodoUpdateRequest struct {
	Status     uint      `json:"status"`
	UpdatedAt  time.Time `json:"updatedAt"`
	UpdateName string    `json:"updateName"`
}

func QueryTodoByOwner(c fiber.Ctx, db *gorm.DB) error {
	// URL decoding
	owner, err := url.QueryUnescape(c.Params("owner"))
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusBadRequest, "客戶端資料錯誤", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	responseData, err := model.FindTodosByOwner(db, owner)
	if err != nil {
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			logrus.Error(err)
			response := ResponseFactory[any](c, fiber.StatusNotFound, err.Error(), nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	logrus.Infof("查詢%s的Todo資料成功", owner)
	response := ResponseFactory(c, fiber.StatusOK, fmt.Sprintf("查詢%s的Todo資料成功", owner), &responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateTodo(c fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData model.Todo
	if err := c.Bind().Body(&clientData); err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusBadRequest, "客戶端資料錯誤", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// handle topic
	lastSlashIndex := strings.LastIndex(clientData.Topic, "/")
	if lastSlashIndex != -1 {
		clientData.Topic = clientData.Topic[:lastSlashIndex]
	} else {
		logrus.Error("topic value has error")
		response := ResponseFactory[any](c, fiber.StatusBadRequest, "客戶端資料轉換錯誤", nil)
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

	err := model.CreateTodo(db, &data)
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	responseData, err := getTodo(c, db)
	if err != nil {
		logrus.Error(err)
		if err == gorm.ErrRecordNotFound {
			response := ResponseFactory[any](c, fiber.StatusNotFound, "找不到該Todo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	logrus.Infof("資料 %s 創建成功", clientData.Title)
	response := ResponseFactory(c, fiber.StatusOK, fmt.Sprintf("資料 %s 創建成功", clientData.Title), responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func UpdateTodoStatus(c fiber.Ctx, db *gorm.DB) error {

	// load client data
	var clientData TodoUpdateRequest
	if err := c.Bind().Body(&clientData); err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusBadRequest, "客戶端資料錯誤", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	clientData.UpdatedAt = time.Now()
	err := model.UpdateTodoStatus(db, c.Params("id"), &model.Todo{Status: clientData.Status, UpdatedAt: clientData.UpdatedAt, UpdateName: clientData.UpdateName})
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := ResponseFactory[any](c, fiber.StatusNotFound, "找不到該Todo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	responseData, err := getTodo(c, db)
	if err != nil {
		logrus.Error(err)
		if err == gorm.ErrRecordNotFound {
			response := ResponseFactory[any](c, fiber.StatusNotFound, "找不到該Todo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	logrus.Infof("Todo %s 更新成功", c.Params("id"))
	response := ResponseFactory(c, fiber.StatusOK, fmt.Sprintf("Todo %s 更新成功", c.Params("id")), responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteTodo(c fiber.Ctx, db *gorm.DB) error {
	err := model.DeleteTodo(db, c.Params("id"))
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := ResponseFactory[any](c, fiber.StatusNotFound, "找不到該Todo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	responseData, err := getTodo(c, db)
	if err != nil {
		logrus.Error(err)
		if err == gorm.ErrRecordNotFound {
			response := ResponseFactory[any](c, fiber.StatusNotFound, "找不到該Todo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	logrus.Infof("Todo %s 刪除成功", c.Params("id"))
	response := ResponseFactory(c, fiber.StatusOK, fmt.Sprintf("Todo %s 刪除成功", c.Params("id")), responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

// 用使用者登入資料取得該使用者的全部Todo資料
// 用在增、改、刪三個API的回傳值，降低前端Request的次數
func getTodo(c fiber.Ctx, db *gorm.DB) (*[]model.Todo, error) {
	userInfo, ok := c.Locals("user_info").(*UserInfo)
	if !ok {
		logrus.Fatal("使用者登入版號表異常")
	}
	responseData, err := model.FindTodosByOwner(db, userInfo.Username)
	if err != nil {
		return nil, err
	}
	return &responseData, nil
}
