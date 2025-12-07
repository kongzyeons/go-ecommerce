package handler

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/pkg/hub"
	"app-ecommerce/pkg/web"
	"bufio"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type notificationHandler struct {
	webSse hub.WebSse
}

func NewNotificationHandler() web.HTTPHandler {
	return &notificationHandler{
		webSse: hub.NewWebSse(),
	}
}

func (h *notificationHandler) Init(router fiber.Router) {
	router.Get("/sse", h.Event)
}

// SseEvent godoc
// @Summary      SSE Event
// @Description  SSE Event
// @Tags         SSE
// @Accept       json
// @Produce      text/event-stream
// @Success      200  {string}  string	"ok"
// @Security     ApiKeyAuth
// @Router       /api/sse [get]
func (h *notificationHandler) Event(c *fiber.Ctx) error {
	userInfo, _ := c.Locals("userInfo").(data.AuthUserInfo)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// use nginx as front-end
	// this is part of trial-error hell to make it work
	c.Set("X-Accel-Buffering", "no")

	client := h.webSse.Register(userInfo.UserID, userInfo.Role)
	// Create a channel for client disconnection
	notify := c.Context().Done()

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		hub.SendMessage(client, w, notify)
	}))
	return nil
}
