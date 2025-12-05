package session

import (
	"app-ecommerce/config"
	"app-ecommerce/internal/meta"
	"sync"

	"github.com/gofiber/fiber/v2/middleware/session"
)

var sesionStore *session.Store
var once sync.Once

func InitSessionStorage() {
	once.Do(func() {
		cfg := config.GetConfig()

		options := session.Config{
			KeyLookup:      "cookie:" + cfg.Web.CookieSessionKey,
			CookieDomain:   cfg.Web.CookieDomain,
			CookieSameSite: cfg.Web.SameSite,
			CookieSecure:   cfg.Web.Secure,
			CookieHTTPOnly: cfg.Web.HTTPOnly,
			// CookieSessionOnly: true,
			Expiration: meta.TTL_AUTH,
		}

		if cfg.ServerMode != "local" {
			options.CookieSecure = true
		}

		if cfg.ServerMode != "local" && cfg.ServerMode != "dev" && cfg.ServerMode != "staging" {
			options.CookieSameSite = "None"
		} else {
			options.CookieSameSite = "None"
		}

		if cfg.Web.UseRedisSession {
			options.Storage = NewFiberRedis()
		}

		sesionStore = session.New(options)

	})
}
