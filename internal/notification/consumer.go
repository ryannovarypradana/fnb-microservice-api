package notification

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Consumer interface {
	StartConsuming()
}

type rabbitMQConsumer struct {
	conn *amqp.Connection
	url  string
}

func NewRabbitMQConsumer(url string) (Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &rabbitMQConsumer{conn: conn, url: url}, nil
}

func (c *rabbitMQConsumer) StartConsuming() {
	ch, err := c.conn.Channel()
	failOnError(err, "Gagal membuka channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"fnb_events", // nama exchange
		"topic",      // tipe exchange
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Gagal mendeklarasikan exchange")

	q, err := ch.QueueDeclare(
		"notification_queue", // nama queue
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	failOnError(err, "Gagal mendeklarasikan queue")

	// Binding queue ke exchange untuk menerima semua event user (user.*)
	err = ch.QueueBind(
		q.Name,       // queue name
		"user.*",     // routing key
		"fnb_events", // exchange
		false,
		nil,
	)
	failOnError(err, "Gagal melakukan binding queue untuk user.*")

	// Binding queue yang SAMA untuk menerima semua event order (order.*)
	err = ch.QueueBind(
		q.Name,       // queue name
		"order.*",    // routing key
		"fnb_events", // exchange
		false,
		nil,
	)
	failOnError(err, "Gagal melakukan binding queue untuk order.*")

	msgs, err := ch.Consume(
		q.Name, "", true, false, false, false, nil,
	)
	failOnError(err, "Gagal mendaftarkan consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Menerima event dengan routing key: %s", d.RoutingKey)
			switch d.RoutingKey {
			case "user.registered":
				handleUserRegistered(d.Body)
			case "order.created":
				handleOrderCreated(d.Body)
			default:
				log.Printf("Warning: Tidak ada handler untuk routing key '%s'", d.RoutingKey)
			}
		}
	}()

	log.Printf(" [*] Notification service menunggu pesan. Untuk keluar tekan CTRL+C")
	<-forever
}

func handleUserRegistered(body []byte) {
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error unmarshalling event body: %s", err)
		return
	}

	email, _ := payload["email"].(string)
	name, _ := payload["name"].(string)

	log.Printf("SIMULASI: Mengirimkan email ke %s...", email)
	time.Sleep(2 * time.Second)
	log.Printf("SUCCESS: Email selamat datang telah dikirim ke %s (Nama: %s)", email, name)
}

func handleOrderCreated(body []byte) {
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error unmarshalling event body: %s", err)
		return
	}

	orderID, _ := payload["order_id"].(string)

	log.Printf("NOTIFIKASI DAPUR: Pesanan baru diterima! Order ID: %s. Mulai siapkan.", orderID)
	time.Sleep(1 * time.Second)
	log.Printf("SUCCESS: Notifikasi untuk Order ID %s telah dikirim ke dapur.", orderID)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
