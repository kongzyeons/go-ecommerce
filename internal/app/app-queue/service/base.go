package service

import (
	"app-ecommerce/internal/events"
	"errors"
	"fmt"
	"log"
)

type EventSvc interface {
	Handle(topic string, eventBytes []byte) error
	Register(req map[string]events.EventHandler)
}

type evnetSvc struct {
	eventHandler map[string]events.EventHandler
}

func NewEventSvc() EventSvc {
	return &evnetSvc{}
}

func (svc *evnetSvc) Handle(topic string, eventBytes []byte) error {
	if _, ok := svc.eventHandler[topic]; !ok {
		return errors.New("topic not found")
	}
	if err := svc.eventHandler[topic].Handle(topic, eventBytes); err != nil {
		return err
	}
	log.Println(fmt.Sprintf("success handle event %s", topic))
	return nil
}

func (svc *evnetSvc) Register(req map[string]events.EventHandler) {
	svc.eventHandler = req
}
