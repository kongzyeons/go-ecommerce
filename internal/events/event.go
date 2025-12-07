package events

import "reflect"

var Topics = []string{
	reflect.TypeOf(OrderEvent{}).Name(),
}

type Event interface {
}

type EventHandler interface {
	Handle(topic string, eventBytes []byte) error
}
