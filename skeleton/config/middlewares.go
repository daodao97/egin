package config

import (
	"github.com/gin-gonic/gin"

	"github.com/daodao97/egin/middleware"
)

type MiddlewaresSlice []func() gin.HandlerFunc

// 由上而下顺序执行
var HttpMiddlewares = MiddlewaresSlice{
	middleware.Cors,
	middleware.MacAuth,
	middleware.IPAuth,
	middleware.IpLimiter,
	middleware.HttpLog,
	middleware.Prometheus,
}
