package events

import (
	"context"

	"github.com/tech-nimble/go-tools/events/queue"
	"github.com/tech-nimble/go-tools/helpers/jaeger"
)

type Repository interface {
	Create(ctx context.Context, entity *Event) error
	Update(ctx context.Context, entity *Event) error
}

type EventBus struct {
	client     *queue.RabbitMQ
	repository Repository
}

func NewEventBus(client *queue.RabbitMQ, repository Repository) *EventBus {
	return &EventBus{
		client:     client,
		repository: repository,
	}
}

func (e *EventBus) Send(ctx context.Context, event *Event) error {
	span, _ := jaeger.StartSpanFromContext(ctx, "EventBus::Send")
	defer span.Finish()

	if err := jaeger.InjectSpanContextToAmqp(span, event); err != nil {
		return err
	}

	if err := e.client.Publish(event); err != nil {
		return err
	}

	event.Sent()
	return e.repository.Update(ctx, event)
}
