package lock

import (
	"context"
	"fmt"
	goredis "github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"myGin/redis"
	"sync"
	"time"
)

type lock struct {
	key         string
	expiration  time.Duration
	requestId   string
	checkCancel chan bool
	mu          sync.Mutex
}

func NewLock(key string, expiration time.Duration) *lock {

	requestId := uuid.NewV4().String()

	return &lock{key: key, expiration: expiration, requestId: requestId, mu: sync.Mutex{}}
}

// Get 获取锁
func (lk *lock) Get() bool {

	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	ok, err := redis.Client().SetNX(cxt, lk.key, lk.requestId, lk.expiration).Result()

	if err != nil {

		return false
	}

	if ok {

		//锁续期检查
		go lk.checkLockIsRelease()
	}

	return ok
}

//检查锁是否被释放，未被释放就延长锁时间
func (lk *lock) checkLockIsRelease() {

	for {

		checkCxt, _ := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(lk.expiration.Milliseconds()-lk.expiration.Milliseconds()/10))

		lk.checkCancel = make(chan bool)

		select {

		case <-checkCxt.Done():

			//多次续期，直到锁被释放
			isContinue := lk.done()

			if !isContinue {

				return
			}

		//取消
		case <-lk.checkCancel:

			fmt.Println("释放")

			return

		}

	}
}

//判断锁是否已被释放
func (lk *lock) done() bool {

	lk.mu.Lock()

	defer lk.mu.Unlock()

	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	res, err := redis.Client().Exists(cxt, lk.key).Result()

	cancel()

	if err != nil {

		return false
	}

	if res == 1 {

		cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		ok, err := redis.Client().Expire(cxt, lk.key, lk.expiration).Result()

		cancel()

		if err != nil {

			return false
		}

		if ok {

			fmt.Println("续期")

			return true

		}

	}

	return false
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

			go lk.checkLockIsRelease()

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

	lk.mu.Lock()

	defer lk.mu.Unlock()

	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := redis.Client().Del(cxt, lk.key).Result()

	if lk.checkCancel != nil {

		lk.checkCancel <- true
	}

	return err

}

// Release 释放锁
func (lk *lock) Release() error {

	lk.mu.Lock()

	defer lk.mu.Unlock()

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

	res, err := script.Run(cxt, redis.Client(), []string{lk.key}, lk.requestId).Result()

	if res.(int64) != 0 {

		lk.checkCancel <- true
	}

	return err

}
