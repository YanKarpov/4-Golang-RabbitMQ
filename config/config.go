package config

import "rabbitmq/redis"  // Подключаем вашу обёртку Redis

type Config struct {
    RabbitMQ struct {
        Login    string
        Password string
        Host     string
        Port     string
    }
    Redis *redis.RedisClient // Используем вашу обёртку RedisClient
}

// NewConfig создаёт новую конфигурацию для RabbitMQ и Redis
func NewConfig(rabbitLogin, rabbitPassword, rabbitHost, rabbitPort, redisAddr, redisPassword string, redisDB int) *Config {
    redisClient := redis.NewRedisClient(redisAddr, redisPassword, redisDB) // Создаём экземпляр RedisClient

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
        Redis: redisClient, // Сохраняем обёртку RedisClient
    }
}
