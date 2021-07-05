package limiter

//component/limiter/limiter.go

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type Limiters struct {
	limiters map[string]*Limiter
	lock     sync.Mutex
}

type Limiter struct {
	limiter *rate.Limiter
	lastGet time.Time //上一次获取token的时间
	key     string
}

var GlobalLimiters = &Limiters{
	limiters: make(map[string]*Limiter),
	lock:     sync.Mutex{},
}

var once = sync.Once{}

func NewLimiter(r rate.Limit, b int, key string) *Limiter {

	once.Do(func() {

		go GlobalLimiters.clearLimiter()
	})

	keyLimiter := GlobalLimiters.getLimiter(r, b, key)

	return keyLimiter

}

func (l *Limiter) Allow() bool {

	l.lastGet = time.Now()

	return l.limiter.Allow()

}

func (ls *Limiters) getLimiter(r rate.Limit, b int, key string) *Limiter {

	ls.lock.Lock()

	defer ls.lock.Unlock()

	limiter, ok := ls.limiters[key]

	if ok {

		return limiter
	}

	l := &Limiter{
		limiter: rate.NewLimiter(r, b),
		lastGet: time.Now(),
		key:     key,
	}

	ls.limiters[key] = l

	return l
}

//清除过期的限流器
func (ls *Limiters) clearLimiter() {

	for {

		time.Sleep(1 * time.Minute)

		for i, i2 := range ls.limiters {

			//超过1分钟
			if time.Now().Unix()-i2.lastGet.Unix() > 60 {
				ls.lock.Lock()
				delete(ls.limiters, i)
				ls.lock.Unlock()
			}

		}

	}

}
