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
	redisClient := redis.NewRedisClient("redis:6379", "", 0) 


	rabbitMQConfig := config.NewConfig("guest", "guest", "rabbitmq", "5672", "/", "5s", 10)

	api := api.NewAPI(redisClient)
	consumer := consumer.NewConsumer(rabbitMQConfig) 

	r := gin.Default()
	r.GET("/api/v1/stats", api.GetStats)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.Consume(ctx) 

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
