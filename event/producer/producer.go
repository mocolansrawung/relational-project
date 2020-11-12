package producer

import (
	"github.com/evermos/boilerplate-go/event/model"
)

type Producer interface {
	Send(event model.EventWrapper, channel string)
}
