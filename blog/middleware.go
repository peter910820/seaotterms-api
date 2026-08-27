package blog

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
)

type UserInfo struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Exp         int       `json:"exp"`
	Management  bool      `json:"management"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdateName  string    `json:"update_name"`
	Avatar      string    `json:"avatar"`
	DataVersion int       `json:"dataVersion"`
}

var (
	// 建立共用使用者快取表，用來處理同使用者登入不同瀏覽器的狀況
	UserInfoInstance = map[uint]*UserInfo{}
)

// 用Token檢查使用者資料(預設前端全域註冊、回傳)
func GetUserInfo(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			logrus.Fatal(err)
		}
		userID := sess.Get("id")
		if userID == nil {
			return c.Next()
		}
		userInfo, ok := UserInfoInstance[userID.(uint)]
		if !ok {
			logrus.Fatal(err) // 有Session但維護表遺失
		}
		c.Locals("user_info", userInfo)
		return c.Next()
	}
}

// 檢查有無登入，沒登入直接中止API並回傳錯誤
func CheckLogin(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userInfo, err := checkLogin(c, store)
		if err != nil {
			logrus.Warn(err)
			response := ResponseFactory[any](c, fiber.StatusUnauthorized, err.Error(), nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}
		c.Locals("user_info", userInfo)
		return c.Next()
	}
}

// 檢查是不是網站管理者，沒登入直接中止API並回傳錯誤
func CheckManagement(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userInfo, err := checkLogin(c, store)
		if err != nil {
			logrus.Warn(err)
			response := ResponseFactory[any](c, fiber.StatusUnauthorized, err.Error(), nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}
		if !userInfo.Management {
			logrus.Warnf("使用者 %s 權限不足", userInfo.Username)
			response := ResponseFactory[any](c, fiber.StatusForbidden, "使用者沒有權限", nil)
			return c.Status(fiber.StatusForbidden).JSON(response)
		}
		c.Locals("user_info", userInfo)
		return c.Next()
	}
}

// utils
func checkLogin(c *fiber.Ctx, store *session.Store) (*UserInfo, error) {
	sess, err := store.Get(c)
	if err != nil {
		logrus.Fatal(err)
	}
	userID := sess.Get("id")
	if userID == nil {
		return nil, errors.New("使用者未登入")
	}

	userInfo, ok := UserInfoInstance[userID.(uint)]
	if !ok {
		return nil, errors.New("使用者未登入")
	}

	return userInfo, nil
}
