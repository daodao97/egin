package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/daodao97/egin/egin/cache/redis"
	"github.com/daodao97/egin/egin/utils/config"
	"github.com/daodao97/egin/egin/utils/limiter"
)

// 单个Api在一秒内的请求次数限制, 不区分用户
func ApiLimiter(limiter *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter.Incr()
		if limiter.CheckOverLimit() {
			c.String(http.StatusBadGateway, "reject")
			c.Abort()
		}
		c.Next()
	}
}

// 单个客户端IP在一秒内的请求次数限制, 不区分Api
func IpLimiter() gin.HandlerFunc {
	red := redis.NewDefault()
	return func(c *gin.Context) {
		conf := config.Config.Auth.IpLimiter
		if !conf.Enable {
			return
		}
		ip := c.ClientIP()
		mu := sync.Mutex{}
		mu.Lock()
		limitCount, ok := conf.IPLimit[ip]
		mu.Unlock()
		if !ok {
			return
		}
		key := fmt.Sprintf("%s:%s", "egin_ip_limiter", c.ClientIP())
		currentCount, _ := red.Get(key)
		_currentCount, _ := strconv.Atoi(currentCount)
		// FIXME 由于 incr 在后, 所以会比实际limit多一次
		// incr放在前又会每次请求都透传到redis, 综合考虑选择后置
		if _currentCount > limitCount {
			c.String(http.StatusBadGateway, "reject")
			c.Abort()
			return
		}
		_ = red.Incr(key)
		_ = red.PExpire(key, time.Second)
		c.Next()
	}
}
