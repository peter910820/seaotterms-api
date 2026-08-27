package blog

import (
	"fmt"
	"net/url"
	model "seaotterms-db/blog"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Request
type (
	UserUpdateRequest struct {
		ID         uint      `json:"id"`
		Username   string    `json:"username"`
		Email      string    `json:"email"`
		Exp        int       `json:"exp"`
		Management bool      `json:"management"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		UpdateName string    `json:"update_name"`
		Avatar     string    `json:"avatar"`
	}

	RegisterRequest struct {
		Username      string `json:"username"`
		Email         string `json:"email"`
		Password      string `json:"password"`
		CheckPassword string `json:"checkPassword"`
	}
)

type UserQueryResponse struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Username   string    `gorm:"NOT NULL unique" json:"username"`
	Exp        int       `gorm:"default:0" json:"exp"`
	Management bool      `gorm:"default:false" json:"management"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	UpdateName string    `json:"updateName"`
}

type UserDataForUpdate struct {
	UpdatedAt  time.Time
	UpdateName string
	Management bool
	Avatar     string
}

type apiAccount struct {
	Username string
	Email    string
}

func QueryUser(c *fiber.Ctx, db *gorm.DB) error {
	var data []UserQueryResponse
	err := db.Table("users").Order("COALESCE(updated_at, created_at) DESC").Find(&data).Error
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := ResponseFactory(c, fiber.StatusOK, "使用者查詢成功", &data)
	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateUser(c *fiber.Ctx, db *gorm.DB) error {
	var data RegisterRequest
	var find []apiAccount

	if err := c.BodyParser(&data); err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusBadRequest, "客戶端資料錯誤", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err := db.Model(&model.User{}).Find(&find).Error
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusOK, fmt.Sprintf("Todo %s 刪除成功", c.Params("id")), nil)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	data.Username = strings.ToLower(data.Username)
	data.Email = strings.ToLower(data.Email)
	// check Username & Email exist
	for _, col := range find {
		if data.Username == col.Username {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, "用戶已註冊", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		} else if data.Email == col.Email {
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, "電子信箱已註冊", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	data.Password = string(hashedPassword)
	dataCreate := model.User{
		Username:   data.Username,
		Password:   data.Password,
		Email:      data.Email,
		CreateName: data.Username,
	}
	err = db.Create(&dataCreate).Error
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := ResponseFactory[any](c, fiber.StatusOK, "註冊成功", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

func UpdateUser(c *fiber.Ctx, db *gorm.DB, store *session.Store) error {
	// load client data
	var clientData UserUpdateRequest
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// URL decoding
	id, err := url.QueryUnescape(c.Params("id"))
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// check if form id equal route id
	u, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logrus.Error(err)
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	if u != uint64(clientData.ID) {
		logrus.Error("ID比對失敗")
		response := ResponseFactory[any](c, fiber.StatusBadRequest, "ID比對失敗", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userInfo, ok := c.Locals("user_info").(*UserInfo)
	if !ok {
		logrus.Error("Middleware異常")
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, "Middleware異常", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	// 如果要更新root，但使用者並不是root，就阻擋更新
	if clientData.Username == "root" && userInfo.Username != "root" {
		logrus.Error("不允許更新root使用者")
		response := ResponseFactory[any](c, fiber.StatusForbidden, "不允許更新root使用者", nil)
		return c.Status(fiber.StatusForbidden).JSON(response)
	}

	timeNow := time.Now()
	if uint(u) == userInfo.ID {
		err = db.Model(&model.User{}).Where("id = ?", u).
			Select("updated_at", "update_name", "avatar").
			Updates(UserDataForUpdate{
				UpdatedAt:  timeNow,
				UpdateName: clientData.Username,
				Avatar:     clientData.Avatar,
			}).Error
	} else {
		err = db.Model(&model.User{}).Where("id = ?", u).
			Select("updated_at", "update_name", "management").
			Updates(UserDataForUpdate{
				UpdatedAt:  timeNow,
				UpdateName: clientData.Username,
				Management: clientData.Management,
			}).Error
	}
	if err != nil {
		// if record not exist
		if err == gorm.ErrRecordNotFound {
			logrus.Error(err)
			response := ResponseFactory[any](c, fiber.StatusNotFound, "使用者不存在，更新使用者失敗", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		} else {
			logrus.Error(err)
			response := ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// 有可能更新的是別人的資料，所以這邊要判斷是更新誰的，更新別人的就不更新快取表
	if userInfo.ID != clientData.ID {
		userInfoCache, ok := UserInfoInstance[clientData.ID]
		// 該使用者可能沒有登入，有登入就會有快取表，更新快取表
		if ok {
			userInfoCache.UpdatedAt = timeNow
			userInfoCache.UpdateName = clientData.Username
			userInfoCache.Management = clientData.Management
			userInfoCache.DataVersion++
		}
		return QueryUser(c, db) // 直接重查(偷懶)
	}

	// 更新使用者資料成功，更新快取表
	userInfoCache, ok := UserInfoInstance[clientData.ID]
	if !ok {
		logrus.Error("快取表更新失敗")
		response := ResponseFactory[any](c, fiber.StatusInternalServerError, "快取表更新失敗", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	userInfoCache.UpdatedAt = timeNow
	userInfoCache.UpdateName = clientData.Username
	userInfoCache.Avatar = clientData.Avatar
	userInfoCache.DataVersion++

	c.Locals("user_info", userInfoCache)

	logrus.Infof("個人資料 %s 更新成功", clientData.Username)
	response := ResponseFactory[any](c, fiber.StatusOK, fmt.Sprintf("資料 %s 更新成功", clientData.Username), nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
