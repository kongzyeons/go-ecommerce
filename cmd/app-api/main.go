package main

import (
	"app-ecommerce/config"
	appapi "app-ecommerce/internal/app/app-api"
	appqueue "app-ecommerce/internal/app/app-queue"
	"encoding/json"
	"fmt"
	"sync"
)

func main() {
	cfg := config.GetConfig()
	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		panic(fmt.Errorf("json marshal config failed, %v", err))
	}
	fmt.Printf("Config: %+s\n", string(jsonCfg))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		appqueue.Run()
	}()

	go func() {
		wg.Wait()
	}()

	appapi.Run()
}
