package main

import (
    "log"
    "context"
    "rabbitmq/api"
    "rabbitmq/config"
    "rabbitmq/consumer"
    "rabbitmq/redis"
    "rabbitmq/producer"
    "github.com/gin-gonic/gin"
)

func main() {
    log.Println("Инициализация клиента Redis")
    redisClient := redis.NewRedisClient("redis:6379", "", 0)
    
    // Проверка подключения к Redis
    if err := redisClient.Ping(context.Background()); err != nil {
        log.Fatalf("Не удалось подключиться к Redis: %v", err)
    } else {
        log.Println("Подключение к Redis успешно!")
    }

    log.Println("Инициализация конфигурации RabbitMQ")
    rabbitMQConfig := config.NewConfig("guest", "guest", "rabbitmq", "5672", "/", "5s", 10)

    log.Println("Создание Producer")
    producer := producer.NewProducer(rabbitMQConfig)

    log.Println("Создание API")
    api := api.NewAPI(redisClient, producer)

    log.Println("Создание Consumer с RedisClient")
    consumer := consumer.NewConsumer(rabbitMQConfig, redisClient)  // Передаем redisClient сюда

    log.Println("Инициализация маршрутов Gin")
    r := gin.Default()

    // Обновленные маршруты
    r.GET("/api/v1/receive/messages/hello", api.GetReceivedMessagesCountByType)
    r.GET("/api/v1/sent/messages/hello", api.GetSentMessagesCountByType)
    r.POST("/api/v1/send/message/hello", api.SendMessageHello)
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
