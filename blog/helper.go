package blog

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func CheckUserAgent(c fiber.Ctx) bool {
	userAgent := c.Get("User-Agent")
	if userAgent == "" || !strings.Contains(userAgent, "Mozilla") {
		return false
	}
	return true
}
