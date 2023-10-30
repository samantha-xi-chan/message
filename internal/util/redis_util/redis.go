package redisutil

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(address, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{client}, nil
}

func (rc *RedisClient) SetWithExpiration(key, value string, expiration time.Duration) error {
	err := rc.client.Set(key, value, expiration).Err()
	return err
}

func (rc *RedisClient) UpdateExpiration(key string, expiration time.Duration) (bool, error) {
	updated, err := rc.client.Expire(key, expiration).Result()
	if err != nil {
		return false, err
	}
	return updated, nil
}

func (rc *RedisClient) Close() {
	rc.client.Close()
}
