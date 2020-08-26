package main

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

type DbInterface interface {
	Save(key string, value string) error
	AddToList(list string, value string) error
	Retrieve(key string) (string, error)
	ReadFromList(list string, start int64, end int64) ([]string, error)
}

// Redis implementation of DbInterface
type RedisDb struct {
	client *redis.Client
}

func NewRedisDb() (*RedisDb, error) {
	// Create Redis Client
	client := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	_, err := client.Ping().Result()
	return &RedisDb{client: client}, err
}

func (db *RedisDb) Save(key string, value string) error {
	// the redis is not a persistent one! If it is, the expiration can be set to 0
	exp := 24 * time.Hour // 0
	cmd, err := db.client.Set(key, value, exp).Result()
	log.Debug(cmd)
	return err
}

func (db *RedisDb) Retrieve(key string) (string, error) {
	return db.client.Get(key).Result()
}

func (db *RedisDb) AddToList(list string, value string) error {
	// doing a push at head
	cmd, err := db.client.LPush(list, value).Result()
	log.Debug(cmd)
	return err
}

func (db *RedisDb) ReadFromList(list string, start int64, end int64) ([]string, error) {
	return db.client.LRange(list, start, end).Result()
}
