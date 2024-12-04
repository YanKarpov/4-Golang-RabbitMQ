package redis

import (
    "context"
    "github.com/redis/go-redis/v9"
)

type RedisClient struct {
    rdb *redis.Client
}

func NewRedisClient(address, password string, db int) *RedisClient {
    return &RedisClient{
        rdb: redis.NewClient(
            &redis.Options{
                Addr:     address,
                Password: password,
                DB:       db,
            },
        ),
    }
}

func (c *RedisClient) Get(key string) (string, error) {
    return c.rdb.Get(context.Background(), key).Result()
}

func (c *RedisClient) Increment(key string) error {
    _, err := c.rdb.Incr(context.Background(), key).Result()
    return err
}

func (c *RedisClient) GetReceivedMessagesCount() (int64, error) {
    return c.rdb.Get(context.Background(), "received_messages").Int64()
}

func (c *RedisClient) GetSentMessagesCount() (int64, error) {
    return c.rdb.Get(context.Background(), "sent_messages").Int64()
}
