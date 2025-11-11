package blog

import (
	"github.com/gofiber/fiber/v2"

	dto "seaotterms-api/dto/blog"
)

func ResponseFactory[T any](c *fiber.Ctx, httpStatus int, msg string, data *T) dto.CommonResponse[T] {
	response := dto.CommonResponse[T]{}
	userInfo, ok := c.Locals("user_info").(*dto.UserInfo)
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
