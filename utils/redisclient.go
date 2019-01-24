package utils

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"github.com/panzhenyu12/flower/config"
)

var rlock sync.Mutex

var redisclient *redis.Client

func GetRedisClient() *redis.Client {
	addr := config.GetConfig().RedisAddr
	if redisclient == nil {
		rlock.Lock()
		defer rlock.Unlock()
		if redisclient == nil {
			glog.Info("init redisclient")
			redisclient = redis.NewClient(&redis.Options{
				Network:     "tcp",
				Addr:        addr,
				PoolSize:    5,
				MaxRetries:  3,
				IdleTimeout: 5 * time.Minute,
			})
		}
	}
	return redisclient
}
