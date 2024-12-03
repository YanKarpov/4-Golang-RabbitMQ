package consumer

import (
	"log"
	"sync"
	 amqp"github.com/rabbitmq/amqp091-go"
	"rabbitmq/config"
)

type Consumer struct {
	Config *config.Config
}

func NewConsumer(config *config.Config) *Consumer {
	return &Consumer{Config: config}
}

// Consume подключается к RabbitMQ и обрабатывает сообщения параллельно
func (c *Consumer) Consume() {
	// Подключение к RabbitMQ
	conn, err := amqp.Dial("amqp://" + c.Config.Login + ":" + c.Config.Password + "@" + c.Config.Host + ":" + c.Config.Port + "/")
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Создание канала
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось открыть канал: %s", err)
	}
	defer ch.Close()

	// Объявление очереди
	q, err := ch.QueueDeclare(
		"YanK", // имя очереди
		true,   // устойчивая
		false,  // автоудаление
		false,  // эксклюзивная
		false,  // без ожидания
		nil,    // аргументы
	)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %s", err)
	}

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
	}

	// Используем sync.WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// Чтение сообщений
	log.Println("Ожидание сообщений. Для выхода нажмите CTRL+C")
	for d := range msgs {
		// Запускаем горутину для обработки каждого сообщения
		wg.Add(1)
		go func(d amqp.Delivery) {
			defer wg.Done()

			// Обработка сообщения через handler
			err := HandleMessage(d)
			if err == nil {
				d.Ack(false) // Подтверждаем обработку
			} else {
				d.Nack(false, true) // Возвращаем сообщение в очередь
			}
		}(d)
	}

	// Ожидаем завершения всех горутин
	wg.Wait()
}

