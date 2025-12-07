package appqueue

import (
	"app-ecommerce/internal/app/app-queue/route"
	"app-ecommerce/internal/events"
	"app-ecommerce/pkg/kafka"
	"context"
	"fmt"
	"os"
	"os/signal"
)

func Run() {
	defer func() {
		kafka.CloseConsumer()
	}()
	Init()

	// Create context that cancels on OS interrupt
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app := route.InitRoute()
	consumer := kafka.GetConsumer()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down gracefully...")
			return
		default:
			consumer.RunConsumer(ctx, events.Topics, app)
		}
	}

}

func Init() {
	// init kafka consumer
	kafka.InitConsumer()
}
