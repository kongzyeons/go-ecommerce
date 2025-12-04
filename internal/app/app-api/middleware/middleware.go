package middleware

import "github.com/gofiber/fiber/v2"

type Middleware interface {
	Authorization(c *fiber.Ctx) error
}

type middleware struct {
}

func NewMiddleware() Middleware {
	return &middleware{}
}

func (m *middleware) Authorization(c *fiber.Ctx) error {
	return c.Next()
}
