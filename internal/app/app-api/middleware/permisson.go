package middleware

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func CheckRoles(roles map[string]bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userInfo, ok := c.Locals("userInfo").(data.AuthUserInfo)
		if !ok {
			return response.Unauthorized[any]().JSON(c)
		}
		if _, ok := roles[userInfo.Role]; !ok {
			return response.Unauthorized[any]().JSON(c)
		}
		return c.Next()
	}
}
