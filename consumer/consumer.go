package consumer

import (
	"context"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"rabbitmq/config"
)

type Consumer struct {
	Config *config.Config
}

func NewConsumer(config *config.Config) *Consumer {
	return &Consumer{Config: config}
}

// Consume подключается к RabbitMQ и обрабатывает сообщения параллельно
func (c *Consumer) Consume(ctx context.Context) {
	var conn *amqp.Connection
	var err error

	// Повторная попытка подключения
	for i := 1; i <= 3; i++ {
		conn, err = amqp.Dial("amqp://" + c.Config.RabbitMQ.Login + ":" + c.Config.RabbitMQ.Password + "@" + c.Config.RabbitMQ.Host + ":" + c.Config.RabbitMQ.Port + "/")
		if err == nil {
			break
		}
		log.Printf("Попытка %d подключения к RabbitMQ не удалась: %s", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ после 3 попыток: %s", err)
		return
	}
	defer conn.Close()
	log.Println("Успешное подключение к RabbitMQ.")

	// Создание канала
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось открыть канал RabbitMQ: %s", err)
		return
	}
	defer ch.Close()

	// Здесь не нужно повторно объявлять очередь. Просто подписываемся на существующую очередь
	q, err := ch.QueueDeclare(
		"MyQueue", // имя очереди (то же, что и в Producer)
		true,      // устойчивая
		false,     // автоудаление
		false,     // эксклюзивная
		false,     // без ожидания
		nil,       // аргументы
	)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %s", err)
		return
	}
	log.Printf("Очередь '%s' успешно объявлена.", q.Name)

	// Подписка на очередь
	msgs, err := ch.Consume(
		q.Name, // имя очереди
		"",     // имя консюмера
		false,  // отключаем авто-подтверждение
		false,  // эксклюзивная
		false,  // без ожидания
		false,  // no-local
		nil,    // аргументы
	)
	if err != nil {
		log.Fatalf("Не удалось зарегистрировать консюмера: %s", err)
		return
	}

	// Ограничение количества горутин
	const maxWorkers = 10
	semaphore := make(chan struct{}, maxWorkers)

	var wg sync.WaitGroup
	log.Println("Ожидание сообщений. Для выхода нажмите CTRL+C")

	for {
		select {
		case <-ctx.Done():
			log.Println("Получен сигнал завершения, останавливаем консюмера.")
			wg.Wait()
			return
		case d, ok := <-msgs:
			if !ok {
				log.Println("Канал сообщений закрыт. Ожидание завершения горутин.")
				wg.Wait()
				return
			}

			semaphore <- struct{}{} // Блокируем до выполнения горутины
			wg.Add(1)

			go func(d amqp.Delivery) {
				defer wg.Done()
				defer func() { <-semaphore }() // Освобождаем место для другой горутины

				// Обработка сообщения
				err := HandleMessage(d, c.Config.Redis)
				if err == nil {
					// Подтверждение успешной обработки сообщения
					if ackErr := d.Ack(false); ackErr != nil {
						log.Printf("Ошибка при подтверждении сообщения: %s", ackErr)
					}
				} else {
					// Возврат сообщения в очередь в случае ошибки
					if nackErr := d.Nack(false, true); nackErr != nil {
						log.Printf("Ошибка при возврате сообщения в очередь: %s", nackErr)
					}
				}
			}(d)
		}
	}
}
