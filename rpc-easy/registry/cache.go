package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type Cache interface {
	GetItem(string) any
	SetItem(string, any)
}

// CacheManager 缓存管理器
var (
	client *redis.Client
	once   sync.Once
)

// RedisCache redis实现本地注册中心
type redisCache struct {
	cli *redis.Client
}

// 避免重复连接redis
func getRedisCli(addr, password string, db int) {
	once.Do(func() {
		if client == nil {
			client = redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: password, // no password set
				DB:       db,       // use default DB
			})
		}
	})
}

func NewRedisCache(addr, password string, db int) Cache {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB

	})
	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Println("redis connect error:", err)
		return nil
	}
	return &redisCache{
		cli: client,
	}
}

func GetRedisCache() Cache {
	getRedisCli("localhost:6379", "", 0)
	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Println("redis connect error:", err)
		return nil
	}
	return &redisCache{
		cli: client,
	}
}

func (r *redisCache) GetItem(key string) any {
	ctx, stop := context.WithTimeout(context.Background(), 3*time.Second)
	val, err := r.cli.Get(ctx, key).Result()
	defer stop()
	if err != nil {
		fmt.Println("redis get error:", err)
		return nil
	}
	return val
}

func (r *redisCache) SetItem(key string, val any) {
	ctx, stop := context.WithTimeout(context.Background(), 3*time.Second)
	defer stop()
	_val, err := json.Marshal(val)
	if err != nil {
		fmt.Println("json marshal error:", err)
	}
	_, err = r.cli.Set(ctx, key, _val, 60*time.Second).Result()

	if err != nil {
		fmt.Println("redis set error:", err)
	}
}

var _mp localCache // 本地缓存, 单例

type localCache struct {
	mp sync.Map
}

func GetLocalCache() Cache {
	return &_mp
}

func (l *localCache) SetItem(serviceName string, handler any) {
	l.mp.Store(serviceName, handler)
}

func (l *localCache) GetItem(serviceName string) any {
	val, ok := l.mp.Load(serviceName)
	if !ok {
		fmt.Println("service not found")
		return nil
	}
	return val
}
