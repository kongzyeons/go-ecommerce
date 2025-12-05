package handler

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/internal/app/app-api/middleware"
	"app-ecommerce/internal/app/app-api/service"
	"app-ecommerce/pkg/response"
	"app-ecommerce/pkg/session"
	"app-ecommerce/pkg/web"
	"time"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	authSvc service.AuthSvc
	session session.SessionStorage
	mw      middleware.Middleware
}

func NewAuthHandler() web.HTTPHandler {
	return &authHandler{
		authSvc: service.NewAuthSvc(),
		session: session.NewSessionStorage(),
		mw:      middleware.NewMiddleware(),
	}
}

func (h *authHandler) Init(router fiber.Router) {
	router.Post("/auth/register", h.Register)
	router.Post("/auth/login", h.Login)
	router.Get("/auth/logout", h.Logout)
	router.Get("/auth/refresh", h.mw.Authorization, h.Refresh)
	router.Get("/auth/user-info", h.mw.Authorization, h.UserInfo)
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

// AuthLogin godoc
// @summary auth login
// @description auth login
// @tags auth
// @security ApiKeyAuth
// @id AuthLogin
// @accept json
// @produce json
// @param AuthLoginReq body data.AuthLoginReq true "request body"
// @success 200 "OK"
// @Router /api/auth/login [post]
func (h *authHandler) Login(c *fiber.Ctx) error {
	var req data.AuthLoginReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest[any]().JSON(c)
	}
	res := h.authSvc.Login(c.Context(), req)
	if res.Success {
		if err := h.session.Set(c, res.Data.UserID, res.Data); err != nil {
			return response.InternalServerError[any](err, "error set session").JSON(c)
		}
	}
	return res.JSON(c)
}

// AuthLogout godoc
// @summary auth logout
// @description auth logout
// @tags auth
// @security ApiKeyAuth
// @id AuthLogout
// @accept json
// @produce json
// @Router /api/auth/logout [get]
func (h *authHandler) Logout(c *fiber.Ctx) error {
	h.session.Delete(c)
	return response.Ok[any](nil).JSON(c)
}

// AuthRefresh godoc
// @summary auth refresh
// @description auth refresh
// @tags auth
// @security ApiKeyAuth
// @id AuthRefresh
// @accept json
// @produce json
// @Router /api/auth/refresh [get]
func (h *authHandler) Refresh(c *fiber.Ctx) error {
	userInfo, _ := c.Locals("userInfo").(data.AuthUserInfo)
	userInfo.LastPing = time.Now().UTC().Unix()
	if err := h.session.Set(c, userInfo.UserID, userInfo); err != nil {
		return response.InternalServerError[any](err, "error set session").JSON(c)
	}
	return response.Ok[any](nil).JSON(c)
}

// AuthUserInfo godoc
// @summary auth user info
// @description auth user info
// @tags auth
// @security ApiKeyAuth
// @id AuthUserInfo
// @accept json
// @produce json
// @success 200 {object} data.AuthUserInfo "OK"
// @Router /api/auth/user-info [get]
func (h *authHandler) UserInfo(c *fiber.Ctx) error {
	userInfo, _ := c.Locals("userInfo").(data.AuthUserInfo)
	return response.Ok(userInfo).JSON(c)
}
