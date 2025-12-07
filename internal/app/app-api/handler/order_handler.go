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
	adminRole := middleware.CheckRoles(map[string]bool{
		"admin": true,
	})

	router.Post("/order/create", userRole, h.Create)
	router.Post("/order/history", userRole, h.GetHistory)
	router.Post("/admin/order/history", adminRole, h.GetHistoryAdmin)
	router.Delete("/order/:id", userRole, h.Delete)
	router.Put("/order/confirm/:id", userRole, h.Confirm)
	router.Put("/order/shipping/:id", adminRole, h.Shipping)
	router.Put("/order/completed/:id", adminRole, h.Completed)
	router.Put("/order/cancel/:id", adminRole, h.Cancel)
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

// OrderGetHistoryAdmin godoc
// @summary admin order get history
// @description admin order get history
// @tags order
// @security ApiKeyAuth
// @id OrderGetHistoryAdmin
// @accept json
// @produce json
// @param OrderGetHistoryReq body data.OrderGetHistoryReq true "request body"
// @success 200 {object} data.OrderGetHistoryRes "OK"
// @Router /api/admin/order/history [post]
func (h *orderHandler) GetHistoryAdmin(c *fiber.Ctx) error {
	var req data.OrderGetHistoryReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
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

// OrderConfirm godoc
// @summary order confirm
// @description orderConfirm
// @tags order
// @security ApiKeyAuth
// @id OrderConfirm
// @accept json
// @produce json
// @param id path int true "order ID"
// @success 200 {object} data.OrderConfirmRes "OK"
// @Router /api/order/confirm/{id} [put]
func (h *orderHandler) Confirm(c *fiber.Ctx) error {
	userInfo, _ := c.Locals("userInfo").(data.AuthUserInfo)
	var req data.OrderConfirmReq
	if err := c.ParamsParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
	req.UserID = userInfo.UserID
	req.ModifiedBy = userInfo.UserName
	res := h.orderSvc.Confirm(c.Context(), req)
	return res.JSON(c)
}

// OrderShipping godoc
// @summary order Shipping
// @description orderShipping
// @tags order
// @security ApiKeyAuth
// @id OrderShipping
// @accept json
// @produce json
// @param id path int true "order ID"
// @success 200 {object} data.OrderShippingRes "OK"
// @Router /api/order/shipping/{id} [put]
func (h *orderHandler) Shipping(c *fiber.Ctx) error {
	userInfo, _ := c.Locals("userInfo").(data.AuthUserInfo)
	var req data.OrderShippingReq
	if err := c.ParamsParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
	req.ModifiedBy = userInfo.UserName
	res := h.orderSvc.Shipping(c.Context(), req)
	return res.JSON(c)
}

// OrderCompleted godoc
// @summary order Completed
// @description orderCompleted
// @tags order
// @security ApiKeyAuth
// @id OrderCompleted
// @accept json
// @produce json
// @param id path int true "order ID"
// @success 200 {object} data.OrderCompletedRes "OK"
// @Router /api/order/completed/{id} [put]
func (h *orderHandler) Completed(c *fiber.Ctx) error {
	userInfo, _ := c.Locals("userInfo").(data.AuthUserInfo)
	var req data.OrderCompletedReq
	if err := c.ParamsParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
	req.ModifiedBy = userInfo.UserName
	res := h.orderSvc.Completed(c.Context(), req)
	return res.JSON(c)
}

// OrderCancel godoc
// @summary order Cancel
// @description orderCancel
// @tags order
// @security ApiKeyAuth
// @id OrderCancel
// @accept json
// @produce json
// @param id path int true "order ID"
// @param OrderCancelReq body data.OrderCancelReq true "request body"
// @success 200 {object} data.OrderCancelRes "OK"
// @Router /api/order/cancel/{id} [put]
func (h *orderHandler) Cancel(c *fiber.Ctx) error {
	userInfo, _ := c.Locals("userInfo").(data.AuthUserInfo)
	var req data.OrderCancelReq
	if err := c.ParamsParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
	req.ModifiedBy = userInfo.UserName
	res := h.orderSvc.Cancel(c.Context(), req)
	return res.JSON(c)
}
