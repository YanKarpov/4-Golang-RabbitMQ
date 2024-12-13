package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	amqp "github.com/rabbitmq/amqp091-go"
)

// HandleMessage обрабатывает входящие сообщения, возможно используя Redis
func HandleMessage(d amqp.Delivery, redisClient *redis.Client) error {
	log.Printf("Получено сообщение: %s", string(d.Body))

	// Проверка подключения к Redis
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Ошибка подключения к Redis: %s", err)
		return err
	}

	

	// Пример использования Redis для инкремента счетчика
	err = redisClient.Incr(context.Background(), "message_count").Err()
	if err != nil {
		log.Printf("Ошибка при инкрементировании счетчика сообщений: %s", err)
		return err
	}

	// Логируем заголовки сообщения для диагностики
	log.Printf("Заголовки сообщения: %+v", d.Headers)

	// Проверяем заголовок type
	messageType, ok := d.Headers["type"].(string)
	if ok && messageType == "hello" {
		log.Println("Сообщение типа hello обработано.")
		return nil
	}

	log.Println("Неизвестный тип сообщения.")
	return fmt.Errorf("неизвестный тип сообщения")
}
