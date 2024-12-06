package producer

import (
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
	"rabbitmq/config"
)

type Producer struct {
	Config *config.Config
}

func NewProducer(config *config.Config) *Producer {
	return &Producer{Config: config}
}

func (p *Producer) Publish(message string) {
	// Подключение к RabbitMQ
	conn, err := amqp.Dial("amqp://" + p.Config.RabbitMQ.Login + ":" + p.Config.RabbitMQ.Password + "@" + p.Config.RabbitMQ.Host + ":" + p.Config.RabbitMQ.Port + "/")
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
		"my_queue", // имя очереди
		true,       // устойчивая
		false,      // автоудаление
		false,      // эксклюзивная
		false,      // без ожидания
		nil,        // аргументы
	)
	if err != nil {
		log.Fatalf("Не удалось объявить очередь: %s", err)
	}

	// Отправка сообщения
	err = ch.Publish(
		"",          // обменник
		q.Name,      // имя очереди
		false,       // обязательность доставки
		false,       // временная очередь
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Fatalf("Не удалось отправить сообщение: %s", err)
	}

	log.Printf("Отправлено сообщение: %s", message)
}