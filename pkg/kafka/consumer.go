package kafka

import (
	"app-ecommerce/internal/events"

	"github.com/IBM/sarama"
)

type consumer struct {
	eventHandler events.EventHandler
}

func NewConsumer(eventHandler events.EventHandler) sarama.ConsumerGroupHandler {
	return consumer{
		eventHandler: eventHandler,
	}
}

func (obj consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (obj consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (obj consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		obj.eventHandler.Handle(msg.Topic, msg.Value)
		session.MarkMessage(msg, "")
	}

	return nil
}
