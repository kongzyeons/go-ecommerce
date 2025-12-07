package kafka

import (
	"app-ecommerce/internal/events"
	"encoding/json"
	"reflect"

	"github.com/IBM/sarama"
)

type Producer interface {
	Send(event events.Event) error
}

type producer struct {
}

func NewProducer() Producer {
	return producer{}
}

func (p producer) Send(event events.Event) error {
	topic := reflect.TypeOf(event).Name()

	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	_, _, err = producerInsatnce.SendMessage(&msg)
	if err != nil {
		return err
	}

	return nil
}
