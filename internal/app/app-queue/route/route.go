package route

import (
	"app-ecommerce/internal/app/app-queue/service"
	"app-ecommerce/internal/events"
	"app-ecommerce/pkg/kafka"
	"reflect"

	"github.com/IBM/sarama"
)

func InitRoute() sarama.ConsumerGroupHandler {
	evnetSvc := service.NewEventSvc()
	evnetSvc.Register(map[string]events.EventHandler{
		reflect.TypeOf(events.OrderEvent{}).Name(): service.NewOrderSvc(),
	})

	consumer := kafka.NewConsumer(evnetSvc)
	return consumer
}
