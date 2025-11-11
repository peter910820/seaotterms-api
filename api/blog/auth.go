package blog

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	dto "seaotterms-api/dto/blog"
	middleware "seaotterms-api/middleware/blog"
	model "seaotterms-api/model/blog"
	utils "seaotterms-api/utils/blog"
)

// 取得使用者資料(本質上是會先進行GetUserInfo這個middlware)
func Auth(c *fiber.Ctx, store *session.Store) error {
	response := utils.ResponseFactory[any](c, fiber.StatusOK, "取得使用者資料成功", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}

// login api
func Login(c *fiber.Ctx, store *session.Store, db *gorm.DB) error {
	var data dto.LoginRequest
	var databaseData []dto.LoginRequest

	if err := c.BodyParser(&data); err != nil {
		logrus.Error(err)
		response := utils.ResponseFactory[any](c, fiber.StatusBadRequest, err.Error(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err := db.Model(&model.User{}).Find(&databaseData).Error
	if err != nil {
		response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	data.Username = strings.ToLower(data.Username)
	for _, col := range databaseData {
		// 找到使用者
		if data.Username == col.Username {
			err := bcrypt.CompareHashAndPassword([]byte(col.Password), []byte(data.Password))
			if err != nil {
				if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
					logrus.Error("login error: password not correct")
					response := utils.ResponseFactory[any](c, fiber.StatusUnauthorized, "密碼輸入錯誤", nil)
					return c.Status(fiber.StatusUnauthorized).JSON(response)
				} else {
					logrus.Error(err)
					response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
					return c.Status(fiber.StatusInternalServerError).JSON(response)
				}
			}

			var userData model.User

			err = db.Where("username = ?", data.Username).First(&userData).Error
			if err != nil {
				logrus.Error(err)
				response := utils.ResponseFactory[any](c, fiber.StatusInternalServerError, err.Error(), nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			data := dto.UserInfo{
				ID:          userData.ID,
				Username:    userData.Username,
				Email:       userData.Email,
				Exp:         userData.Exp,
				Management:  userData.Management,
				CreatedAt:   userData.CreatedAt,
				UpdatedAt:   userData.UpdatedAt,
				UpdateName:  userData.UpdateName,
				Avatar:      userData.Avatar,
				DataVersion: 1,
			}

			// 如果有登入紀錄，因為有重查一次DB，所以更新一次資料以及版號
			value, ok := middleware.UserInfo[userData.ID]
			if ok {
				value.Username = data.Username
				value.Email = data.Email
				value.Exp = data.Exp
				value.Management = data.Management
				value.CreatedAt = data.CreatedAt
				value.UpdatedAt = data.UpdatedAt
				value.UpdateName = data.UpdateName
				data.Avatar = userData.Avatar
				value.DataVersion++
			} else {
				middleware.UserInfo[userData.ID] = &data
			}

			c.Locals("user_info", middleware.UserInfo[userData.ID])

			setUserInfoSession(c, store, &data)

			response := utils.ResponseFactory[any](c, fiber.StatusOK, fmt.Sprintf("使用者 %s 登入成功", data.Username), nil)
			return c.Status(fiber.StatusOK).JSON(response)
		}
	}
	logrus.Error("user not found")
	response := utils.ResponseFactory[any](c, fiber.StatusUnauthorized, "找不到該使用者: "+data.Username, nil)
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}

func setUserInfoSession(c *fiber.Ctx, store *session.Store, userInfo *dto.UserInfo) {
	// set session
	sess, err := store.Get(c)
	if err != nil {
		logrus.Fatal(err) // 這邊之後會發送訊息
	}

	sess.Set("id", userInfo.ID)
	if err := sess.Save(); err != nil {
		logrus.Fatal(err) // 這邊之後會發送訊息
	}

	logrus.Infof("Username %s login success", userInfo.Username)
}
