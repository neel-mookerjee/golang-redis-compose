package main

import (
	"github.com/go-redis/redis"
)

type EndpointRepoInterface interface {
	GetEndpoints() ([]string, error)
}

// Redis implementation of EndpointRepoInterface: NOT COOL bcz it points to the same dataspurce lol
type EndpointRepo struct {
	client *redis.Client
}

func NewEndpointRepo() (*EndpointRepo, error) {
	// Create Redis Client
	client := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	_, err := client.Ping().Result()
	return &EndpointRepo{client: client}, err
}

// get all health endpoints
func (e *EndpointRepo) GetEndpoints() ([]string, error) {
	keys, err := e.client.Keys("url:*").Result()
	if err != nil {
		return nil, err
	}
	op := keys
	for i, key := range keys {
		url, err := e.client.Get(key).Result()
		if err != nil {
			return nil, err
		}
		op[i] = url + "/health"
	}
	return op, nil

}
