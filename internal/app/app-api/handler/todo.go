package handler

import (
	"app-ecommerce/pkg/response"
	"app-ecommerce/pkg/web"

	"github.com/gofiber/fiber/v2"
)

type todo struct {
}

func NewTodo() web.HTTPHandler {
	return &todo{}
}

func (h *todo) Init(router fiber.Router) {
	router.Get("/todo/ping", h.Ping)
}

// TodoPing godoc
// @summary todo ping
// @description todo ping
// @tags todo
// @security ApiKeyAuth
// @id TodoPing
// @accept json
// @produce json
// @success 200 {object} string "OK"
// @Router /api/todo/ping [get]
func (h *todo) Ping(c *fiber.Ctx) error {
	res := response.Ok("pong")
	return res.JSON(c)
}
