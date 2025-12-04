package handler

import (
	"app-ecommerce/pkg/web"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
)

type swaggerHandler struct {
}

func NewSwaggerHandler() web.HTTPHandler {
	return swaggerHandler{}
}

func (h swaggerHandler) Init(router fiber.Router) {
	// isAdmin := middleware.CheckRoles(map[string]bool{"admin": true})
	router.Get("/docs", h.Document)
}

func (h swaggerHandler) Document(c *fiber.Ctx) error {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./internal/app/app-api/docs/swagger.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Exta API",
		},
		DarkMode: true,
	})

	if err != nil {
		return err
	}
	c.Type("html")
	return c.SendString(htmlContent)
}
