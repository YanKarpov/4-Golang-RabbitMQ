package config

import "github.com/redis/go-redis/v9"

type Config struct {
    RabbitMQ struct {
        Login    string
        Password string
        Host     string
        Port     string
    }
    Redis *redis.Client // Используем *redis.Client
}

// NewConfig создаёт новую конфигурацию для RabbitMQ и Redis
func NewConfig(rabbitLogin, rabbitPassword, rabbitHost, rabbitPort, redisAddr, redisPassword string, redisDB int) *Config {
    redisClient := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: redisPassword,
        DB:       redisDB,
    })

    return &Config{
        RabbitMQ: struct {
            Login    string
            Password string
            Host     string
            Port     string
        }{
            Login:    rabbitLogin,
            Password: rabbitPassword,
            Host:     rabbitHost,
            Port:     rabbitPort,
        },
        Redis: redisClient, // Сохраняем *redis.Client
    }
}
