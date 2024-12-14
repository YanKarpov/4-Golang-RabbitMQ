package consumer

import (
	"context"
	"log"
	"rabbitmq/redis" // Импортируем ваш пакет с обёрткой
	amqp "github.com/rabbitmq/amqp091-go"
)

// HandleMessage обрабатывает входящие сообщения и записывает их в Redis
func HandleMessage(d amqp.Delivery, redisClient *redis.RedisClient) error {
	// Логируем сообщение
	log.Printf("Получено сообщение: %s", string(d.Body))

	// Логируем заголовки сообщения для диагностики
	log.Printf("Заголовки сообщения: %+v", d.Headers)

	// Проверяем подключение к Redis с использованием обёртки
	err := redisClient.Ping(context.Background())
	if err != nil {
		log.Printf("Ошибка при подключении к Redis: %s", err)
		return err
	}

	// Записываем сообщение в Redis (например, добавляем в список)
	err = redisClient.LPush(context.Background(), "received_messages", string(d.Body)) 
	if err != nil {
		log.Printf("Ошибка при записи в Redis: %s", err)
		return err
	}

	// Логируем, что сообщение успешно записано в Redis
	log.Println("Сообщение успешно записано в Redis.")

	return nil
}
