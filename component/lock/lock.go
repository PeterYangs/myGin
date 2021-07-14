package lock

import (
	"context"
	goredis "github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"myGin/redis"
	"time"
)

type lock struct {
	key        string
	expiration time.Duration
	requestId  string
}

func NewLock(key string, expiration time.Duration) *lock {

	requestId := uuid.NewV4().String()

	return &lock{key: key, expiration: expiration, requestId: requestId}
}

// Get 获取锁
func (lk *lock) Get() bool {

	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	ok, err := redis.Client().SetNX(cxt, lk.key, lk.requestId, lk.expiration).Result()

	if err != nil {

		return false
	}

	return ok
}

// Block 阻塞获取锁
func (lk *lock) Block(expiration time.Duration) bool {

	t := time.Now()

	for {

		cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		ok, err := redis.Client().SetNX(cxt, lk.key, lk.requestId, lk.expiration).Result()

		cancel()

		if err != nil {

			return false
		}

		if ok {

			return true
		}

		time.Sleep(200 * time.Millisecond)

		if time.Now().Sub(t) > expiration {

			return false
		}

	}

}

// ForceRelease 强制释放锁，忽略请求id
func (lk *lock) ForceRelease() error {

	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := redis.Client().Del(cxt, lk.key).Result()

	return err

}

// Release 释放锁
func (lk *lock) Release() error {

	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	const luaScript = `
	if redis.call('get', KEYS[1])==ARGV[1] then
		return redis.call('del', KEYS[1])
	else
		return 0
	end
	`

	script := goredis.NewScript(luaScript)

	_, err := script.Run(cxt, redis.Client(), []string{lk.key}, lk.requestId).Result()

	return err

}
