package route

import (
	"app-ecommerce/config"
	"app-ecommerce/internal/app/app-api/handler"
	"app-ecommerce/internal/app/app-api/middleware"
	"app-ecommerce/pkg/web"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func InitRoute() *fiber.App {
	app := fiber.New(fiber.Config{EnableSplittingOnParsers: true})
	SetUp(app)

	baseRoot := app.Group("api")

	registerPublicHandlers(baseRoot)

	// middlerware
	mw := middleware.NewMiddleware()
	privateRoot := baseRoot.Use(mw.Authorization)
	registerPrivateHandlers(privateRoot)

	return app
}

func registerPublicHandlers(root fiber.Router) {
	handlers := web.HandlerRegistrator{}
	handlers.Register(
		handler.NewTodo(),
		handler.NewSwaggerHandler(),
		handler.NewAuthHandler(),
	)
	handlers.Init(root)
}

func registerPrivateHandlers(root fiber.Router) {
	handlers := web.HandlerRegistrator{}
	handlers.Register(
		handler.NewProduct(),
		handler.NewOrderHandler(),
	)
	handlers.Init(root)
}

func SetUp(app *fiber.App) {
	cfg := config.GetConfig()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     strings.Join(cfg.Web.CORSAllowOrigin, ","),
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	if cfg.IsDebug {
		app.Use(logger.New(logger.Config{
			Format:     "[${time}] (${pid}) ${status} ${method} ${path} - lat: ${latency} err: ${error} reqid: ${locals:requestid}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   "Asia/Bangkok",
		}))
	}

	// Apply Rate Limiter if it is enabled
	// if config.RateLimiter.RateLimitEnabled {

	// 	rateLimiter := limiter.New(limiter.Config{
	// 		Max:        config.RateLimiter.RateLimitMax,        // default 20
	// 		Expiration: config.RateLimiter.RateLimitExpiration, //default 1 minute
	// 		LimitReached: func(c *fiber.Ctx) error {
	// 			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests, please try again later.")
	// 		},
	// 	})

	// 	// Apply the rate limiter middleware
	// 	app.Use(rateLimiter)
	// }

}
