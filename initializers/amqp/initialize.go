package amqp

import (
	"github.com/gobuffalo/envy"
	"github.com/tech-nimble/go-tools/events"
	"github.com/tech-nimble/go-tools/events/queue"
	"github.com/tech-nimble/go-tools/events/repository"
	"github.com/tech-nimble/go-tools/repositories"
)

const (
	AMQPUrlEnv          = "AMQP_URL"
	AMQPExchangeNameEnv = "AMQP_EXCHANGE_NAME"
	AMQPExchangeTypeEnv = "AMQP_EXCHANGE_TYPE"
)

func InitializeAMQPOption() *queue.Options {
	return &queue.Options{
		Url:             envy.Get(AMQPUrlEnv, ""),
		ConnectAttempts: 5,
		Exchange: struct {
			Name    string
			Type    string
			Durable bool
		}{
			Name:    envy.Get(AMQPExchangeNameEnv, "users"),
			Type:    envy.Get(AMQPExchangeTypeEnv, "topic"),
			Durable: true,
		},
	}
}

func InitializeAMQP(options *queue.Options) (*queue.RabbitMQ, error) {
	client := queue.NewRabbitMQ(options)
	if err := client.Init(); err != nil {
		return nil, err
	}

	return client, nil
}

func InitializeEventRepository(rep *repositories.Repositories) *repository.Events {
	return repository.NewEvents(rep, "events")
}

func InitializeEventBus(client *queue.RabbitMQ, rep *repository.Events) *events.EventBus {
	return events.NewEventBus(client, rep)
}
