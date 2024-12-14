package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

// NewRedisClient создаёт новый RedisClient
func NewRedisClient(address, password string, db int) *RedisClient {
	return &RedisClient{
		rdb: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
	}
}

// Ping выполняет команду PING для проверки подключения
func (c *RedisClient) Ping(ctx context.Context) error {
	_, err := c.rdb.Ping(ctx).Result()
	return err
}

// Get выполняет команду GET для получения значения по ключу
func (c *RedisClient) Get(key string) (string, error) {
	return c.rdb.Get(context.Background(), key).Result()
}

// Increment увеличивает значение ключа на 1
func (c *RedisClient) Increment(key string) error {
	_, err := c.rdb.Incr(context.Background(), key).Result()
	return err
}

// GetReceivedMessagesCount возвращает количество полученных сообщений
func (c *RedisClient) GetReceivedMessagesCount() (int64, error) {
	return c.rdb.Get(context.Background(), "received_messages").Int64()
}

// GetSentMessagesCount возвращает количество отправленных сообщений
func (c *RedisClient) GetSentMessagesCount() (int64, error) {
	return c.rdb.Get(context.Background(), "sent_messages").Int64()
}

// GetMessagesCountByType возвращает количество сообщений по типу и направлению
func (c *RedisClient) GetMessagesCountByType(messageType, direction string) (int64, error) {
	key := messageType + ":" + direction
	return c.rdb.Get(context.Background(), key).Int64()
}

// LPush добавляет элемент в список в Redis
func (c *RedisClient) LPush(ctx context.Context, key, value string) error {
	_, err := c.rdb.LPush(ctx, key, value).Result()
	return err
}

