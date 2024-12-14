package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "rabbitmq/redis"
    "rabbitmq/producer"
)

type API struct {
    rdb       *redis.RedisClient
    producer  *producer.Producer
}

func NewAPI(redisClient *redis.RedisClient, producer *producer.Producer) *API {
    return &API{
        rdb:      redisClient,
        producer: producer,
    }
}

// GetReceivedMessagesCountByType - Получить количество полученных сообщений с типом Hello
func (api *API) GetReceivedMessagesCountByType(c *gin.Context) {
    count, err := api.rdb.GetMessagesCountByType("hello", "received")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить количество полученных сообщений типа Hello"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "received_hello_messages": count,
    })
}

// GetSentMessagesCountByType - Получить количество отправленных сообщений с типом Hello
func (api *API) GetSentMessagesCountByType(c *gin.Context) {
    count, err := api.rdb.GetMessagesCountByType("hello", "sent")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить количество отправленных сообщений типа Hello"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "sent_hello_messages": count,
    })
}

// SendMessageHello - Отправить сообщение типа "Hello" в RabbitMQ
func (api *API) SendMessageHello(c *gin.Context) {
    // Отправляем сообщение через producer
    err := api.producer.Publish("hello")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить сообщение типа Hello"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "ok",
    })
}

func (api *API) GetStats(c *gin.Context) {
    received, err := api.rdb.GetReceivedMessagesCount()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить количество полученных сообщений"})
        return
    }

    sent, err := api.rdb.GetSentMessagesCount()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить количество отправленных сообщений"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "received_messages": received,
        "sent_messages":     sent,
    })
}
