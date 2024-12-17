package consumer

import (
	"context"
	"log"
	"rabbitmq/redis"
	amqp "github.com/rabbitmq/amqp091-go"
)

// HandleMessage обрабатывает входящие сообщения и записывает их в Redis
func HandleMessage(d amqp.Delivery, redisClient *redis.RedisClient) error {
	// Логируем сообщение
	log.Printf("Получено сообщение: %s", string(d.Body))

	// Логируем заголовки сообщения для диагностики
	log.Printf("Заголовки сообщения: %+v", d.Headers)

	// Пытаемся записать сообщение в Redis
	err := redisClient.LPush(context.Background(), "received_messages", string(d.Body))
	if err != nil {
		log.Printf("Ошибка при записи в Redis: %s", err)
		return err
	}

	// Логируем успешную запись в Redis
	log.Println("Сообщение успешно записано в Redis.")
	return nil
}
