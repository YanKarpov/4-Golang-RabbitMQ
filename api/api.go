package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "rabbitmq/redis"
)

type API struct {
    rdb *redis.RedisClient
}

func NewAPI(redisClient *redis.RedisClient) *API {
    return &API{
        rdb: redisClient,
    }
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
