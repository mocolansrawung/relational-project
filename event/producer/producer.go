package producer

type Producer interface {
	Send(event EventWrapper, channel string)
}
