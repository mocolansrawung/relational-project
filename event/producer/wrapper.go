package producer

import (
	"time"
)

type EventWrapper struct {
	EventType string `json:"event_type"`
	Data      Data   `json:"data"`
}

type Data struct {
	Timestamp time.Time   `json:"timestamp"`
	Value     interface{} `json:"value"`
}

func Wrapper(eventType string, model interface{}) EventWrapper {
	return EventWrapper{
		EventType: eventType,
		Data: Data{
			Timestamp: time.Now(),
			Value:     model,
		},
	}
}
