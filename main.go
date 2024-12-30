package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var ctx = context.Background()

type RedisLock struct {
	client *redis.Client
	ttl    time.Duration
	ch     chan struct{}
	key    string
}

func (l *RedisLock) Lock() {
	_, err := l.client.SetNX(ctx, l.key, "", l.ttl).Result()
	if err != nil {
		panic(err)
	}
	log.Println("lock success")
	l.ch = make(chan struct{})
	// 开启守护线程
	go func() {
		select {
		case <-l.ch:
			log.Println("cancel delay")
			return
		default:
			_, err = l.client.Expire(ctx, l.key, 3*time.Second).Result()
			if err != nil {
				panic(err)
			}
			<-time.After(time.Second * 2)
		}
	}()
}

func (l *RedisLock) Unlock() {
	close(l.ch)
	err := l.client.Del(ctx, l.key).Err()
	if err != nil {
		panic(err)
	}
	log.Println("lock release")
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
		ttl:    3 * time.Second,
		key:    "lock",
	}
	buf := &buffer{
		l:    lock,
		buf:  make([]int, 0, 5),
		size: 5,
	}
	go func() {
		for i := 0; i < 4; i++ {
			if i%2 == 0 {
				consume(buf, i)
			} else {
				produce(buf, i)
			}
		}
	}()
	time.Sleep(10 * time.Second)
}

func consume(buf *buffer, id int) {
	buf.l.Lock()
	log.Println("consumer:", id, "lock")
	defer buf.l.Unlock()
	if len(buf.buf) == 0 {
		log.Println("buf is empty")
		return
	}
	item := buf.buf[0]
	buf.buf = buf.buf[1:]
	log.Println("consume item", item)
	log.Println("consumer:", id, "unlock")
}

func produce(buf *buffer, item int) {
	buf.l.Lock()
	log.Println("producer:", item, "lock")
	defer buf.l.Unlock()
	for len(buf.buf) == buf.size {
		log.Println("buf is full")
		return
	}
	buf.buf = append(buf.buf, item)
	log.Println("produce item", item)
	log.Println("producer:", item, "unlock")
}
