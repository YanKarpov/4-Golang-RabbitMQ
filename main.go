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
	log.Println("Инициализация клиента Redis")
	redisClient := redis.NewRedisClient("redis:6379", "", 0)

	log.Println("Инициализация конфигурации RabbitMQ")
	rabbitMQConfig := config.NewConfig("guest", "guest", "rabbitmq", "5672", "/", "5s", 10)

	log.Println("Создание API")
	api := api.NewAPI(redisClient)

	log.Println("Создание Consumer")
	consumer := consumer.NewConsumer(rabbitMQConfig) 

	log.Println("Инициализация маршрутов Gin")
	r := gin.Default()
	r.GET("/api/v1/stats", api.GetStats)

	log.Println("Запуск Consumer в горутине")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.Consume(ctx) 

	log.Println("Запуск HTTP-сервера")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

