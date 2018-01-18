package cache

import (
	"github.com/go-redis/redis"
	"log"
	"sync"
	"time"
	"fmt"
)

type RedisCache struct {
	Client *redis.Client
	CacheExpirationSeconds int
}

var redisCache *RedisCache
var once sync.Once

func BuildInstance(host, port string, cacheExpirationSeconds int) *RedisCache {
	once.Do(func() {
		redisCache = newRedisCache(host, port, cacheExpirationSeconds)
	})
	return redisCache
}

func GetInstance() *RedisCache {
	return redisCache
}

// Builds a new Redis Cache
func newRedisCache(host, port string, cacheExpirationSeconds int) *RedisCache {
	log.Printf("New Cache")

	address := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	client.FlushAll()
	redisCache := RedisCache{
		Client: client,
		CacheExpirationSeconds: cacheExpirationSeconds,
	}

	return &redisCache
}

// get Object from Cache
func (r *RedisCache) Get(key string) (string, error) {
	var err error

	val, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return "", err
	}
	if err != nil {
		panic(err)
	}

	log.Printf("CACHE HIT for key %s", key)
	return val, nil
}

// sets Object in Cache using key
func (r *RedisCache) Set(key, json string) error {
	err := r.Client.Set(key, json, time.Second * time.Duration(r.CacheExpirationSeconds)).Err()
	if err != nil {
		panic(err)
	}
	return err
}
