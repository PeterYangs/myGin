package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type connect struct {
	client *redis.Client
}

var once = sync.Once{}

var _connect *connect

func connectRedis() {

	cxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	conf := &redis.Options{
		Addr: "192.168.232.128:6380",
		DB:   0,
	}

	c := redis.NewClient(conf)

	re := c.Ping(cxt)

	if re.Err() != nil {

		panic(re.Err())

	}

	_connect = &connect{
		client: c,
	}

}

func Client() *redis.Client {

	if _connect == nil {

		once.Do(func() {

			connectRedis()
		})

	}

	return _connect.client

}
