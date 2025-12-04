package main

import (
	"app-ecommerce/config"
	appapi "app-ecommerce/internal/app/app-api"
	"encoding/json"
	"fmt"
)

func main() {
	cfg := config.GetConfig()
	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		panic(fmt.Errorf("json marshal config failed, %v", err))
	}
	fmt.Printf("Config: %+s\n", string(jsonCfg))
	appapi.Run()
}
