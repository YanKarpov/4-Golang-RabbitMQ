package consumer

import (
	"context"
	"log"
	"sync"
	"rabbitmq/redis"
	"rabbitmq/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Config      *config.Config
	RedisClient *redis.RedisClient
}

func NewConsumer(config *config.Config, redisClient *redis.RedisClient) *Consumer {
	return &Consumer{
		Config:      config,
		RedisClient: redisClient,
	}
}

func (c *Consumer) Consume(ctx context.Context) {
	conn, err := amqp.Dial("amqp://" + c.Config.RabbitMQ.Login + ":" + c.Config.RabbitMQ.Password + "@" + c.Config.RabbitMQ.Host + ":" + c.Config.RabbitMQ.Port + "/")
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %s", err)
		return
	}
	defer conn.Close()
	log.Println("Успешное подключение к RabbitMQ.")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось открыть канал RabbitMQ: %s", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"MyQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %s", err)
		return
	}
	log.Printf("Очередь '%s' успешно объявлена.", q.Name)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Не удалось зарегистрировать консюмера: %s", err)
		return
	}

	const maxWorkers = 10
	semaphore := make(chan struct{}, maxWorkers)

	var wg sync.WaitGroup
	log.Println("Ожидание сообщений. Для выхода нажмите CTRL+C")

	if err := c.RedisClient.Ping(ctx); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %s", err)
	}

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

			semaphore <- struct{}{}
			wg.Add(1)

			go func(d amqp.Delivery) {
				defer wg.Done()
				defer func() { <-semaphore }()

				err := HandleMessage(d, c.RedisClient)
				if err == nil {
					if ackErr := d.Ack(false); ackErr != nil {
						log.Printf("Ошибка при подтверждении сообщения: %s", ackErr)
					}
				} else {
					if nackErr := d.Nack(false, true); nackErr != nil {
						log.Printf("Ошибка при возврате сообщения в очередь: %s", nackErr)
					}
				}
			}(d)
		}
	}
}
