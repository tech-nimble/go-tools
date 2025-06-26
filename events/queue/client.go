package queue

import (
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type RabbitMQ struct {
	options         *Options
	connectionError chan *amqp.Error
	connection      *amqp.Connection
	channel         *amqp.Channel
	notifyClose     chan bool
	notifyReConsume chan bool
	isClosed        bool
}

func NewRabbitMQ(options *Options) *RabbitMQ {
	return &RabbitMQ{
		options:         options,
		notifyReConsume: make(chan bool, 1),
		notifyClose:     make(chan bool, 1),
	}
}

func (q *RabbitMQ) Init() error {
	var err error

	if err = q.connect(); err == nil {
		go q.watchReconnect()
	}

	return err
}

func (q *RabbitMQ) GetStatus() bool {
	return !q.isClosed
}

func (q *RabbitMQ) GetExchangeName() string {
	return q.options.Exchange.Name
}

func (q *RabbitMQ) Publish(message Publishing) error {
	body, err := message.GetBody()
	if err != nil {
		return err
	}

	return q.channel.Publish(
		q.options.Exchange.Name,
		message.GetRoutingKey(),
		false,
		false,
		amqp.Publishing{
			Headers:      message.GetHeaders(),
			Body:         body,
			DeliveryMode: uint8(message.GetDeliveryMode()),
		},
	)
}

// закрываем соединение с rabbitmq
func (q *RabbitMQ) Close() {
	q.notifyClose <- true
	defer close(q.notifyClose)
	defer close(q.notifyReConsume)

	q.channel.Close()
}

func (q *RabbitMQ) connect() error {
	var err error

	if q.connection != nil {
		q.connection.Close()
	}

	if q.connection, err = amqp.Dial(q.options.Url); err != nil {
		return err
	}

	if q.channel != nil {
		q.channel.Close()
	}

	if q.channel, err = q.connection.Channel(); err != nil {
		return err
	}

	if err = q.channel.ExchangeDeclare(
		q.options.Exchange.Name,
		q.options.Exchange.Type,
		q.options.Exchange.Durable,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	q.notifyClose = make(chan bool, 1)

	return nil
}

func (q *RabbitMQ) GetChannel() *amqp.Channel {
	return q.channel
}

func (q *RabbitMQ) watchReconnect() {
	q.connectionError = make(chan *amqp.Error)
	q.channel.NotifyClose(q.connectionError)

rcn:
	for {
		select {
		case connErr := <-q.connectionError:
			if err := q.tryReconnect(); err != nil {
				log.Fatal().
					Err(connErr).
					Int("attempt", q.options.ConnectAttempts).
					Msg("reconnect failed")
				q.Close()
				break rcn
			}
		case <-q.notifyClose:
			break rcn
		}
	}
}

func (q *RabbitMQ) tryReconnect() error {
	q.isClosed = true

	for {
		var err error
		if err = q.connect(); err == nil {
			q.connectionError = make(chan *amqp.Error)
			q.channel.NotifyClose(q.connectionError)

			q.isClosed = false
			q.notifyReConsume <- true

			return nil
		}

		log.Error().Err(err).Msg("connect failed")

		time.Sleep(time.Duration(10) * time.Second)
	}

	return errors.New("cannot reconnect to rabbitmq")
}
