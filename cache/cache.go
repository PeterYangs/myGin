package cache

import (
	"myGin/cache/fileCache"
	"myGin/cache/redisCache"
	"time"
)

type CacheContract interface {
	Put(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
}

func Cache(driver string) CacheContract {

	switch driver {

	case "redis":

		return redisCache.NewRedisCache()

	case "file":

		return fileCache.NewFileCache()

	}

	return redisCache.NewRedisCache()
}
