package nunchakus

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type RedisQueue struct {
	queueName string
	length int

	red *redis.Client
}

func NewRedisQueue(opt *redis.Options, prefix string) *RedisQueue {
	if opt == nil {
		opt = &redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}
	}
	rdb := redis.NewClient(opt)
	return &RedisQueue{red: rdb, queueName: prefix + ".queue"}
}

func (rq *RedisQueue) Push(item string) bool {
	_, err := rq.red.LPush(context.Background(), rq.queueName, item).Result()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func (rq *RedisQueue) Pop() string {
	res, err := rq.red.LPop(context.Background(), rq.queueName).Result()
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return res
}

func (rq *RedisQueue) PopN(n int) []string {
	msg, err := rq.red.LPopCount(context.Background(), rq.queueName, n).Result()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return msg
}
