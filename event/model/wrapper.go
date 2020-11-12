package model

import (
	"encoding/json"
	"time"
)

type EventWrapper struct {
	EventType string `json:"event_type"`
	Data      Data   `json:"data"`
}

type Data struct {
	Timestamp time.Time `json:"timestamp"`
	Value     []byte    `json:"value"`
}

func Wrapper(eventType string, model interface{}) EventWrapper {
	value, _ := json.Marshal(model)

	return EventWrapper{
		EventType: eventType,
		Data: Data{
			Timestamp: time.Now(),
			Value:     value,
		},
	}
}
