package consumer

import (
	"log"
	 amqp"github.com/rabbitmq/amqp091-go"
)

// Обработка полученного сообщения
func HandleMessage(d amqp.Delivery) error {
	log.Printf("Обработка сообщения: %s", d.Body)


	return nil
}
