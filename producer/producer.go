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

// Publish отправляет сообщение в RabbitMQ и возвращает ошибку, если она произошла
func (p *Producer) Publish(message string) error {
	conn, err := amqp.Dial("amqp://" + p.Config.RabbitMQ.Login + ":" + p.Config.RabbitMQ.Password + "@" + p.Config.RabbitMQ.Host + ":" + p.Config.RabbitMQ.Port + "/")
	if err != nil {
		log.Printf("Не удалось подключиться к RabbitMQ: %s", err)
		return err // возвращаем ошибку
	}
	defer conn.Close()
	log.Println("Успешное подключение к RabbitMQ.")

	// Создание канала
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Не удалось открыть канал: %s", err)
		return err // возвращаем ошибку
	}
	defer ch.Close()

	// Объявление очереди (Producer и Consumer используют одну очередь)
	q, err := ch.QueueDeclare(
		"MyQueue", // имя очереди
		true,      // устойчивая
		false,     // автоудаление
		false,     // эксклюзивная
		false,     // без ожидания
		nil,       // аргументы
	)
	if err != nil {
		log.Printf("Не удалось объявить очередь: %s", err)
		return err // возвращаем ошибку
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
		log.Printf("Не удалось отправить сообщение: %s", err)
		return err // возвращаем ошибку
	}

	log.Printf("Отправлено сообщение: %s", message)
	return nil // возвращаем nil, если всё прошло успешно
}
