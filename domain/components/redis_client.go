package components

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/noydev/ggmapapi-test/utils/logger"
)

type RedisClient struct {
	redisClient *redis.Client
	timeOut     time.Duration
}

type RedisClientConfig struct {
	Addr     string        `mapstructure:"address"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

func NewRedisClient(cfg RedisClientConfig) RedisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
	return RedisClient{redisClient: c, timeOut: cfg.Timeout}
}

func (rc RedisClient) Set(key string, value interface{}) error {
	valuePacked, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	err = rc.redisClient.Set(context.Background(), key, valuePacked, rc.timeOut).Err()
	if err != nil {
		logger.Error(fmt.Sprintf("Error setting value to Redis: %s", err))
		return err
	}
	return nil
}

func (rc RedisClient) Get(key string, dest interface{}) error {
	val, err := rc.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting value from Redis: %s", err))
		return err
	}
	return json.Unmarshal([]byte(val), &dest)
}
