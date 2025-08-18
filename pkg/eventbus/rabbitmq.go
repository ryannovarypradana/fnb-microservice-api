// pkg/eventbus/rabbitmq.go
package eventbus

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

// EventBus mendefinisikan interface untuk message broker kita.
type EventBus interface {
	Publish(exchange, routingKey string, body interface{}) error
	Close()
}

type rabbitMQBus struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewRabbitMQBus membuat koneksi baru ke RabbitMQ.
func NewRabbitMQBus(url string) (EventBus, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &rabbitMQBus{conn: conn, ch: ch}, nil
}

// Publish mengirim pesan ke exchange dengan routing key tertentu.
func (r *rabbitMQBus) Publish(exchange, routingKey string, body interface{}) error {
	// Konversi body ke JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = r.ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		},
	)
	return err
}

func (r *rabbitMQBus) Close() {
	r.ch.Close()
	r.conn.Close()
}
