package db

import (
	"tinyquant/src/util"

	. "tinyquant/src/logger"

	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client {

	return redisClient
}

func InitRedis() error {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     util.RedisHost,
		Password: util.RedisPass, // no password set
		DB:       3,              // use default DB
		PoolSize: 100,
	})

	if redisClient == nil {
		panic("redis init failed")
	}

	_, err := redisClient.Ping().Result()
	if err != nil {
		Logger.Error("redis connect failed", zap.Error(err))
		return err
	}

	Logger.Info("Redis connect Success")
	return nil
}

func RedisDestory() {
	if redisClient != nil {
		redisClient.Close()
	}
}
