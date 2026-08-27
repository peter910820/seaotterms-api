package blog

import (
	"github.com/gofiber/fiber/v2"
)

type CommonResponse[T any] struct {
	StatusCode int       `json:"statusCode"` // http status code
	ErrMsg     string    `json:"errMsg"`
	InfoMsg    string    `json:"infoMsg"`
	UserInfo   *UserInfo `json:"userInfo"`
	Data       *T        `json:"data"`
}

func ResponseFactory[T any](c *fiber.Ctx, httpStatus int, msg string, data *T) CommonResponse[T] {
	response := CommonResponse[T]{}
	userInfo, ok := c.Locals("user_info").(*UserInfo)
	if ok {
		response.UserInfo = userInfo
		// logrus.Debugf("%s 使用者資料版號: %d", response.UserInfo.Username, response.UserInfo.DataVersion)
	}
	response.StatusCode = httpStatus
	response.Data = data
	if httpStatus == 200 {
		response.InfoMsg = msg
	} else {
		response.ErrMsg = msg
	}
	return response
}
