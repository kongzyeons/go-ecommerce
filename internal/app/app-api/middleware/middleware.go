package middleware

import (
	"app-ecommerce/internal/app/app-api/data"
	"app-ecommerce/pkg/response"
	"app-ecommerce/pkg/session"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Middleware interface {
	Authorization(c *fiber.Ctx) error
}

type middleware struct {
	session session.SessionStorage
}

var middlewareInstance Middleware
var middlewareOnce sync.Once

func NewMiddleware() Middleware {
	middlewareOnce.Do(func() {
		middlewareInstance = &middleware{
			session: session.NewSessionStorage(),
		}
	})
	return middlewareInstance
}

func (m *middleware) Authorization(c *fiber.Ctx) error {
	var userInfo data.AuthUserInfo
	if err := m.session.Get(c, &userInfo); err != nil {
		return response.Unauthorized[any]().JSON(c)
	}
	c.Locals("userInfo", userInfo)
	return c.Next()
}

// func (m *middleware) Authorization(c *fiber.Ctx) error {
// 	cfg := config.GetConfig()
// 	headerAuth := c.Get("Authorization")
// 	token, err := jwt.Parse(headerAuth, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(cfg.JwtKey), nil
// 	})
// 	if err != nil {
// 		return response.Unauthorized[any](err.Error()).JSON(c)
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if ok && token.Valid {
// 		var userInfo data.AuthUserInfo
// 		jsonStr, err := json.Marshal(claims)
// 		if err != nil {
// 			return response.Unauthorized[any](err.Error()).JSON(c)
// 		}

// 		err = json.Unmarshal(jsonStr, &userInfo)
// 		if err != nil {
// 			return response.Unauthorized[any](err.Error()).JSON(c)
// 		}
// 		c.Locals("userInfo", userInfo)
// 	}
// 	return c.Next()
// }
