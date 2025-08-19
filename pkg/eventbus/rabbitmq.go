package eventbus

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// rabbitMQPublisher adalah implementasi konkret dari interface Publisher.
type rabbitMQPublisher struct {
	conn *amqp.Connection
}

// NewRabbitMQPublisher membuat publisher baru dan memastikan ia sesuai dengan interface Publisher.
func NewRabbitMQPublisher(url string) (Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke RabbitMQ: %w", err)
	}
	return &rabbitMQPublisher{conn: conn}, nil
}

func (p *rabbitMQPublisher) Publish(exchange, routingKey, contentType string, body map[string]interface{}) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("gagal membuka channel: %w", err)
	}
	defer ch.Close()

	// Pastikan exchange "topic" ada
	err = ch.ExchangeDeclare(
		exchange, // nama exchange
		"topic",  // tipe
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("gagal mendeklarasikan exchange: %w", err)
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("gagal marshal body ke JSON: %w", err)
	}

	log.Printf("Mengirim pesan ke exchange '%s' dengan routing key '%s'", exchange, routingKey)
	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        bodyBytes,
		},
	)
	if err != nil {
		return fmt.Errorf("gagal mempublikasikan pesan: %w", err)
	}

	return nil
}

func (p *rabbitMQPublisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
