package handler

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/internal/app/app-api/service"
	"app-ecommerce/pkg/response"
	"app-ecommerce/pkg/web"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	authSvc service.AuthSvc
}

func NewAuthHandler() web.HTTPHandler {
	return &authHandler{
		authSvc: service.NewAuthSvc(),
	}
}

func (h *authHandler) Init(router fiber.Router) {
	router.Post("/auth/register", h.Register)
}

// AuthRegister godoc
// @summary auth register
// @description auth register
// @tags auth
// @security ApiKeyAuth
// @id AuthRegister
// @accept json
// @produce json
// @param AuthRegisterReq body data.AuthRegisterReq true "request body"
// @success 200 "OK"
// @Router /api/auth/register [post]
func (h *authHandler) Register(c *fiber.Ctx) error {
	var req data.AuthRegisterReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
	res := h.authSvc.Register(c.Context(), req)
	return res.JSON(c)
}
