package appapi

import (
	"app-ecommerce/config"
	"app-ecommerce/internal/app/app-api/route"
	"app-ecommerce/pkg/db"
	"app-ecommerce/pkg/hub"
	"app-ecommerce/pkg/kafka"
	redis_db "app-ecommerce/pkg/redis"
	"app-ecommerce/pkg/session"
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
// @name Authorization
func Run() {
	defer func() {
		db.UnInitDatabase()
		redis_db.Uninit()
		kafka.CloseProducer()
		hub.CloseHub()
	}()
	Init()

	app := route.InitRoute()
	log.Println("start server...")
	app.Listen(fmt.Sprintf(":%s", config.GetConfig().Web.PORT))
}

func Init() {
	cfg := config.GetConfig()

	// init db
	db.InitDatabase(cfg.PostgresDB.ConnectionString)

	// init redis
	redis_db.Init(
		cfg.Redis.Hosts,
		cfg.Redis.Password,
		cfg.Redis.UseCluster,
		cfg.Redis.UseTLS,
	)

	// init session storage
	session.InitSessionStorage()

	// init producer
	kafka.InitProducer()

	// init hub
	hub.InitHub().Run()
}
