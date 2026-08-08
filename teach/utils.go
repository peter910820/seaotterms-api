package teach

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func dbErrHandler(c *fiber.Ctx, err error) error {
	logrus.Error(err)
	if err == gorm.ErrRecordNotFound {
		//404
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": err.Error(),
		})
	} else {
		// 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
}
