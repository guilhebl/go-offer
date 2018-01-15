package cache

import (
	"github.com/go-redis/redis"
	"log"
	"sync"
	"time"
	"fmt"
	"github.com/guilhebl/go-offer/common/model"
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

// get Object From Cache
func (r *RedisCache) GetOfferList(key string) (*model.OfferList, error) {
	var wanted model.OfferList
	var err error

	val, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return nil, err
	}
	if err != nil {
		panic(err)
	}

	if err := wanted.UnmarshalBinary([]byte(val)); err != nil {
		log.Printf("Unable to unmarshal data from REDIS cache: %s \n", err)
		return nil, err
	}

	log.Printf("CACHE HIT for key %s", key)
	return &wanted, nil
}

// sets Object in Cache using key
func (r *RedisCache) SetOfferList(key string, obj *model.OfferList) error {
	err := r.Client.Set(key, obj, time.Second * time.Duration(r.CacheExpirationSeconds)).Err()
	if err != nil {
		panic(err)
	}
	return err
}
