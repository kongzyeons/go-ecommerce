package handler

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/internal/app/app-api/service"
	"app-ecommerce/pkg/web"

	"github.com/gofiber/fiber/v2"
)

type productHandler struct {
	productSvc service.ProductSvc
}

func NewProduct() web.HTTPHandler {
	return &productHandler{
		productSvc: service.NewProductSvc(),
	}
}

func (h *productHandler) Init(router fiber.Router) {
	router.Post("/products", h.GetList)
}

// ProductGetList godoc
// @summary product get list
// @description product get list
// @tags product
// @security ApiKeyAuth
// @id ProductGetList
// @accept json
// @produce json
// @param ProductGetListReq body data.ProductGetListReq true "request body"
// @success 200 {object} data.ProductGetListRes "OK"
// @Router /api/products [post]
func (h *productHandler) GetList(c *fiber.Ctx) error {
	var req data.ProductGetListReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	res := h.productSvc.GetList(c.Context(), req)
	return res.JSON(c)
}
