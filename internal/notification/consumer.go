// internal/notification/consumer.go
package notification

import (
	"log"

	"github.com/streadway/amqp"
)

func StartOrderCreatedConsumer(rabbitURL string) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Pastikan exchange "orders" ada
	err = ch.ExchangeDeclare("orders", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	// Buat queue sementara yang unik untuk service ini
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Ikat (bind) queue ke exchange dengan routing key "order.created"
	err = ch.QueueBind(q.Name, "order.created", "orders", false, nil)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %v", err)
	}

	// Mulai mengkonsumsi pesan dari queue
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Notification service is waiting for order.created events...")

	// Blok selamanya untuk terus mendengarkan
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received an order.created event: %s", d.Body)
			// Di sini Anda bisa menambahkan logika untuk mengirim email, push notification, dll.
		}
	}()
	<-forever
}
