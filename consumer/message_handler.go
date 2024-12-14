package consumer

import (
	"context"
	"log"
	"rabbitmq/redis" // Импортируем ваш пакет с обёрткой
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

// HandleMessage обрабатывает входящие сообщения и записывает их в Redis
func HandleMessage(d amqp.Delivery, redisClient *redis.RedisClient) error {
	// Логируем сообщение
	log.Printf("Получено сообщение: %s", string(d.Body))

	// Логируем заголовки сообщения для диагностики
	log.Printf("Заголовки сообщения: %+v", d.Headers)

	// Максимальное количество попыток подключения к Redis
	maxRetries := 5
	retryDelay := time.Second // Задержка между попытками (1 секунда)

	// Флаг, показывающий, что подключение к Redis не удалось
	redisError := false

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Проверяем подключение к Redis с использованием обёртки
		err := redisClient.Ping(context.Background())
		if err == nil {
			// Если подключение успешно, продолжаем обработку
			break
		}

		// Если попытки исчерпаны, выводим ошибку
		if attempt == maxRetries {
			log.Printf("Ошибка при подключении к Redis: %s. Превышено количество попыток.", err)
			redisError = true
			break
		}

		// Логируем попытку и ждем перед следующей
		log.Printf("Ошибка при подключении к Redis: %s. Попытка %d из %d. Повтор через %v", err, attempt, maxRetries, retryDelay)
		time.Sleep(retryDelay)
	}

	// Если после всех попыток подключения к Redis ошибка осталась, возвращаем ошибку
	if redisError {
		log.Println("Не удаётся подключиться к Redis, пропускаем сообщение.")
		// Возвращаем ошибку, чтобы не обрабатывать дальнейшие сообщения
		return nil // Или можно использовать ошибку для дальнейшей обработки, например, игнорировать её
	}

	// Записываем сообщение в Redis (например, добавляем в список)
	err := redisClient.LPush(context.Background(), "received_messages", string(d.Body))
	if err != nil {
		log.Printf("Ошибка при записи в Redis: %s", err)
		return err
	}

	// Логируем, что сообщение успешно записано в Redis
	log.Println("Сообщение успешно записано в Redis.")

	// Подтверждаем успешную обработку сообщения
	if err := d.Ack(false); err != nil {
		log.Printf("Ошибка при подтверждении сообщения: %s", err)
		return err
	}

	return nil
}
