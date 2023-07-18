package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient interface {
	Get(string) (string, error)
	Set(string, string) error
	Del([]string) error
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient() RedisClient {
	dbStr := os.Getenv("REDIS_DATABASE")
	db, er := strconv.Atoi(dbStr)
	if er != nil {
		fmt.Println("Can not get redis db")
		db = 0
	}
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
	return &redisClient{
		client: client,
	}
}

func (r *redisClient) Set(key, data string) error {
	err := r.client.Set(ctx, key, data, 0).Err()
	return err
}
func (r *redisClient) Del(key []string) error {
	err := r.client.Del(ctx, key...).Err()
	return err
}

func (r *redisClient) Get(key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", err
	} else if err != nil {
		return "", err
	} else {
		return val, err
	}
}

func (r *redisClient) MSet(key string, data []interface{}) error {
	err := r.client.MSet(ctx, data...).Err()
	return err
}

func (r *redisClient) MGet(key string) ([]interface{}, error) {
	val, err := r.client.MGet(ctx, key).Result()
	if err == redis.Nil {
		return []interface{}{}, err
	} else if err != nil {
		return []interface{}{}, err
	} else {
		return val, err
	}
}
