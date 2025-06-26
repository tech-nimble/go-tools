package events

import (
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/tech-nimble/go-tools/events/queue"
)

const (
	statusNew = iota
	statusSent
)

type EventData interface {
	GetModelID() int
	GetID() string
}

type Event struct {
	ID         int
	Payload    any
	Headers    map[string]any
	Status     int
	ModelID    int
	RoutingKey string
	Exchange   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewEvent(payload any, headers map[string]any, modelID int, routingKey, exchange string) *Event {
	return &Event{
		ID:         0,
		Payload:    payload,
		Headers:    headers,
		Status:     statusNew,
		ModelID:    modelID,
		RoutingKey: routingKey,
		Exchange:   exchange,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Time{},
	}
}

func (e *Event) Sent() {
	e.Status = statusSent
	e.UpdatedAt = time.Now()
}

func (e *Event) GetHeaders() amqp.Table {
	return e.Headers
}

func (e *Event) SetHeaders(headers map[string]any) {
	e.Headers = headers
}

func (e *Event) GetBody() ([]byte, error) {
	return json.Marshal(e.Payload)
}

func (e *Event) GetDeliveryMode() queue.DeliveryMode {
	return queue.DeliveryMode(amqp.Persistent)
}

func (e *Event) GetRoutingKey() string {
	return e.RoutingKey
}
