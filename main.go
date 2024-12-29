package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

type RedisLock struct {
	client *redis.Client
	key    string
	val    string
	ch     chan struct{}
}

func (l RedisLock) Lock() {
	_, err := l.client.SetNX(ctx, l.key, l.val, 3*time.Second).Result()
	if err != nil {
		panic(err)
	}
	// 开启守护线程
	go func() {
		select {
		case <-l.ch:
			return
		default:
			time.Sleep(2 * time.Second)
			l.client.ExpireAt(ctx, l.key, time.Now().Add(3*time.Second))
		}
	}()
}

func (l RedisLock) Unlock() {
	err := l.client.Del(ctx, l.key).Err()
	if err != nil {
		panic(err)
	}
	l.ch <- struct{}{}
}

type buffer struct {
	l    *RedisLock
	buf  []int
	size int
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // use default DB
	})
	lock := &RedisLock{
		client: client,
		key:    "lock",
		val:    "lock",
	}
	buf := &buffer{
		l:    lock,
		buf:  make([]int, 0),
		size: 5,
	}
	go func() {
		for i := 0; i < 5; i++ {
			produce(buf, i)
		}
	}()
	go func() {
		for i := 0; i < 5; i++ {
			consume(buf)
		}
	}()
	time.Sleep(10 * time.Second)
}

func consume(buf *buffer) {
	buf.l.Lock()
	defer buf.l.Unlock()
	if len(buf.buf) == 0 {
		println("buf is empty")
		return
	}
	item := buf.buf[0]
	buf.buf = buf.buf[1:]
	println("consume item", item)
}

func produce(buf *buffer, item int) {
	buf.l.Lock()
	defer buf.l.Unlock()
	for len(buf.buf) == buf.size {
		println("buf is full")
		return
	}
	buf.buf = append(buf.buf, item)
	println("produce item", item)
}
