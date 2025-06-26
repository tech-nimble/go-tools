package queue

import amqp "github.com/rabbitmq/amqp091-go"

const (
	Transient DeliveryMode = 1 + iota
	Persistent
)

type DeliveryMode uint8

type Options struct {
	Url             string
	ConnectAttempts int
	Exchange        struct {
		Name    string
		Type    string
		Durable bool
	}
}

type Publishing interface {
	GetBody() ([]byte, error)
	GetHeaders() amqp.Table
	SetHeaders(map[string]any)
	GetDeliveryMode() DeliveryMode
	GetRoutingKey() string
}
