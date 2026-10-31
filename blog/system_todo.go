package blog

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	model "seaotterms-db/blog"
)

// Request
type (
	SystemTodoCreateRequest struct {
		SystemName  string     `json:"systemName"`
		Title       string     `json:"title"`
		Detail      string     `json:"detail"`
		Status      uint       `json:"status"`
		Deadline    *time.Time `json:"deadline"`
		Urgency     uint       `json:"urgency"`
		CreatedName string     `json:"createdName"`
	}
	SystemTodoUpdateRequest struct {
		SystemName  string     `json:"systemName"`
		Title       string     `json:"title"`
		Detail      string     `json:"detail"`
		Status      uint       `json:"status"`
		Deadline    *time.Time `json:"deadline"`
		Urgency     uint       `json:"urgency"`
		UpdatedAt   time.Time  `json:"updatedAt"`
		UpdatedName string     `json:"updatedName"`
	}

	QuickSystemTodoUpdateRequest struct {
		Status      uint      `json:"status"`
		UpdatedAt   time.Time `json:"updatedAt"`
		UpdatedName string    `json:"updatedName"`
	}
)

func QuerySystemTodo(c fiber.Ctx, db *gorm.DB) error {
	// get query param
	id := c.Query("id")
	systemName := c.Query("system_name")
	statusStr := c.Query("status")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		status = 999
	}

	var data []model.SystemTodo
	switch {
	case id != "":
		data, err = model.FindSystemTodoByID(db, id)
	case systemName != "" && status != 999:
		data, err = model.FindSystemTodosBySystemNameAndStatus(db, systemName, status)
	case systemName != "":
		data, err = model.FindSystemTodosBySystemName(db, systemName)
	case status != 999:
		data, err = model.FindSystemTodosByStatus(db, status)
	default:
		data, err = model.FindAllSystemTodos(db)
	}
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := ResponseFactory[any](c, fiber.StatusNotFound, "找不到SystemTodo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	logrus.Info("查詢SystemTodo資料成功")
	response := ResponseFactory(c, fiber.StatusOK, "查詢SystemTodo資料成功", &data)
	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateSystemTodo(c fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData SystemTodoCreateRequest
	if err := c.Bind().Body(&clientData); err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	data := model.SystemTodo{
		SystemName:  clientData.SystemName,
		Title:       clientData.Title,
		Detail:      clientData.Detail,
		Status:      clientData.Status,
		Deadline:    clientData.Deadline,
		Urgency:     clientData.Urgency,
		CreatedName: clientData.CreatedName,
	}

	err := model.CreateSystemTodo(db, &data)
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	logrus.Infof("系統代辦資料 %s 創建成功", clientData.Title)
	response := ResponseFactory[any](c, fiber.StatusOK, fmt.Sprintf("系統代辦資料 %s 創建成功", clientData.Title), nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

func UpdateSystemTodo(c fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData SystemTodoUpdateRequest
	if err := c.Bind().Body(&clientData); err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	clientData.UpdatedAt = time.Now()
	err := model.UpdateSystemTodo(db, c.Params("id"), &model.SystemTodo{SystemName: clientData.SystemName, Title: clientData.Title, Detail: clientData.Detail, Status: clientData.Status, Deadline: clientData.Deadline, Urgency: clientData.Urgency, UpdatedAt: clientData.UpdatedAt, UpdatedName: clientData.UpdatedName})
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := ResponseFactory[any](c, fiber.StatusNotFound, "找不到該SystemTodo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	logrus.Infof("SystemTodo %s 更新成功", c.Params("id"))
	response := ResponseFactory[any](c, fiber.StatusOK, fmt.Sprintf("SystemTodo %s 更新成功", c.Params("id")), nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

func QuickUpdateSystemTodo(c fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData QuickSystemTodoUpdateRequest
	if err := c.Bind().Body(&clientData); err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	clientData.UpdatedAt = time.Now()
	err := model.QuickUpdateSystemTodo(db, c.Params("id"), &model.SystemTodo{Status: clientData.Status, UpdatedAt: clientData.UpdatedAt, UpdatedName: clientData.UpdatedName})
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := ResponseFactory[any](c, fiber.StatusNotFound, "找不到該SystemTodo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	logrus.Infof("SystemTodo %s 更新成功", c.Params("id"))
	response := ResponseFactory[any](c, fiber.StatusOK, fmt.Sprintf("SystemTodo %s 更新成功", c.Params("id")), nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteSystemTodo(c fiber.Ctx, db *gorm.DB) error {
	err := model.DeleteSystemTodo(db, c.Params("id"))
	if err != nil {
		logrus.Error(err)
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			response := ResponseFactory[any](c, fiber.StatusNotFound, "找不到該SystemTodo資料", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	logrus.Infof("SystemTodo %s 刪除成功", c.Params("id"))
	response := ResponseFactory[any](c, fiber.StatusOK, fmt.Sprintf("SystemTodo %s 刪除成功", c.Params("id")), nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
