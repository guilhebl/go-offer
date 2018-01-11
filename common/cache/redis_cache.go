package cache

import (
	"fmt"
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/guilhebl/go-offer/common/model"
	"github.com/vmihailenco/msgpack"
	"log"
	"sync"
	"time"
)

type RedisCache struct {
	Codec *cache.Codec
}

// gets OfferList From Cache
func (r *RedisCache) GetOfferList(key string) (*model.OfferList, error) {
	var wanted model.OfferList
	if err := r.Codec.Get(key, &wanted); err == nil {
		fmt.Println("ERORR " + err.Error())

		return nil, err
	}
	return &wanted, nil
}

// sets OfferList Into Cache - key JSON representation of request
func (r *RedisCache) SetOfferList(key string, obj *model.OfferList) error {
	err := r.Codec.Set(&cache.Item{
		Key:        key,
		Object:     obj,
		Expiration: time.Hour,
	})
	return err
}

var redisCache *RedisCache
var once sync.Once

func BuildInstance() *RedisCache {
	once.Do(func() {
		redisCache = newRedisCache()
	})
	return redisCache
}

func GetInstance() *RedisCache {
	return redisCache
}

// Builds a new Redis Cache
func newRedisCache() *RedisCache {
	log.Printf("New Cache")

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379",
			"server2": ":6380",
		},
		PoolSize: 10,
	})

	codec := &cache.Codec{
		Redis: ring,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}

	redisCache := RedisCache{
		Codec: codec,
	}

	return &redisCache
}
