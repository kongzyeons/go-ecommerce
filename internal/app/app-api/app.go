package appapi

import (
	"app-ecommerce/config"
	"app-ecommerce/internal/app/app-api/route"
	"fmt"
	"log"
)

// @title  APP Ecommerce API
// @description This is a sample server API.
// @version 1.0
// @BasePath /
// @schemes https http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func Run() {
	defer func() {
	}()

	Init()

	app := route.InitRoute()
	log.Println("start server...")
	app.Listen(fmt.Sprintf(":%s", config.GetConfig().Web.PORT))
}

func Init() {
	config.GetConfig()
}
