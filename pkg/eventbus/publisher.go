package eventbus

// Publisher mendefinisikan interface untuk mempublikasikan event.
// Dengan adanya interface ini, Anda bisa menukar RabbitMQ dengan Kafka atau teknologi lain
// di masa depan tanpa harus mengubah service yang menggunakannya.
type Publisher interface {
	Publish(exchange, routingKey, contentType string, body map[string]interface{}) error
	Close()
}
