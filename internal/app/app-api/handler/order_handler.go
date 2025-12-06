package handler

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/internal/app/app-api/middleware"
	"app-ecommerce/internal/app/app-api/service"
	"app-ecommerce/pkg/response"
	"app-ecommerce/pkg/web"

	"github.com/gofiber/fiber/v2"
)

type orderHandler struct {
	orderSvc service.OrderSvc
}

func NewOrderHandler() web.HTTPHandler {
	return &orderHandler{
		orderSvc: service.NewOrderSvc(),
	}
}

func (h *orderHandler) Init(router fiber.Router) {
	userRole := middleware.CheckRoles(map[string]bool{
		"user": true,
	})

	router.Post("/order/create", userRole, h.Create)
	router.Post("/order/history", userRole, h.GetHistory)
	router.Delete("/order/:id", userRole, h.Delete)
}

// OrderCreate godoc
// @summary order create
// @description order create
// @tags order
// @security ApiKeyAuth
// @id OrderCreate
// @accept json
// @produce json
// @param OrderCreateReq body data.OrderCreateReq true "request body"
// @success 200 {object} data.OrderCreateRes "OK"
// @Router /api/order/create [post]
func (h *orderHandler) Create(c *fiber.Ctx) error {
	userInfo, _ := c.Locals("userInfo").(data.AuthUserInfo)
	var req data.OrderCreateReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
	req.UserID = userInfo.UserID
	req.CreateBy = userInfo.UserName
	res := h.orderSvc.Create(c.Context(), req)
	return res.JSON(c)
}

// OrderGetHistory godoc
// @summary order get history
// @description order get history
// @tags order
// @security ApiKeyAuth
// @id OrderGetHistory
// @accept json
// @produce json
// @param OrderGetHistoryReq body data.OrderGetHistoryReq true "request body"
// @success 200 {object} data.OrderGetHistoryRes "OK"
// @Router /api/order/history [post]
func (h *orderHandler) GetHistory(c *fiber.Ctx) error {
	userInfo, _ := c.Locals("userInfo").(data.AuthUserInfo)
	var req data.OrderGetHistoryReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
	req.UserID = userInfo.UserID
	res := h.orderSvc.GetHistory(c.Context(), req)
	return res.JSON(c)
}

// OrderDelete godoc
// @summary order delete
// @description orderdelete
// @tags order
// @security ApiKeyAuth
// @id OrderDelete
// @accept json
// @produce json
// @param id path int true "order ID"
// @success 200 {object} data.OrderDeleteRes "OK"
// @Router /api/order/{id} [delete]
func (h *orderHandler) Delete(c *fiber.Ctx) error {
	userInfo, _ := c.Locals("userInfo").(data.AuthUserInfo)
	var req data.OrderDeleteReq
	if err := c.ParamsParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
	req.UserID = userInfo.UserID
	req.DeletedBy = userInfo.UserName
	res := h.orderSvc.Delete(c.Context(), req)
	return res.JSON(c)
}
