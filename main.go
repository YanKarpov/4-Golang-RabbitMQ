package main

import (
	"log"
	"rabbitmq/config"
	"rabbitmq/consumer"
	"rabbitmq/producer"
	"sync"
	"strconv"
)

func main() {
	// Настроим конфигурацию RabbitMQ
	config := &config.Config{
		Login:    "guest",   // логин
		Password: "guest",   // пароль
		Host:     "localhost", // хост RabbitMQ
		Port:     "5672",      // порт RabbitMQ
	}

	// Создаем нового потребителя
	consumer := consumer.NewConsumer(config)

	// Создаем производителя
	producer := producer.NewProducer(config)

	var wg sync.WaitGroup

	// Отправим большое количество сообщений
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			msg := "Сообщение номер " + strconv.Itoa(i) // Преобразуем число в строку
			log.Printf("Отправка сообщения: %s", msg)
			producer.Publish(msg)
		}(i)
	}

	// Запускаем консюмера
	go consumer.Consume()

	// Логируем, что консюмер запущен
	log.Println("Консюмер и производитель запущены. Ожидаю сообщения...")

	// Ожидаем завершения всех горутин, связанных с отправкой сообщений
	wg.Wait()
	log.Println("Все сообщения отправлены.")

	// Блокируем основной поток, чтобы не завершить выполнение
	select {}
}
