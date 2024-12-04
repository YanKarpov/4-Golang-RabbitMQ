package main

import (
	"log"
	"context"
	"rabbitmq/api"
	"rabbitmq/config"
	"rabbitmq/consumer"
	"rabbitmq/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализируем Redis клиент
	redisClient := redis.NewRedisClient("localhost:6379", "", 0)

	// Инициализируем объект Config для RabbitMQ
	// Убедитесь, что передаете все необходимые параметры, например, VirtualHost и ReconnectInterval
	rabbitMQConfig := config.NewConfig("guest", "guest", "localhost", "5672", "/", "5s", 10)

	// Инициализируем API и Consumer
	api := api.NewAPI(redisClient)
	consumer := consumer.NewConsumer(rabbitMQConfig) // Передаем rabbitMQConfig, а не redisClient

	// Настроим маршруты API
	r := gin.Default()
	r.GET("/api/v1/stats", api.GetStats)

	// Запускаем Consumer в отдельной горутине с передачей контекста
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.Consume(ctx) // Передаем контекст

	// Запускаем HTTP сервер
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
