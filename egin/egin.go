package egin

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/daodao97/egin/db"
	"github.com/daodao97/egin/egin/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/daodao97/egin/egin/consts"
	"github.com/daodao97/egin/egin/middleware"
	"github.com/daodao97/egin/egin/utils/config"
)

type Bootstrap struct {
	HttpMiddlewares []func() gin.HandlerFunc
	engine          *gin.Engine
	RegRoutes       func(r *gin.Engine)
}

func (boot *Bootstrap) Start() {
	gin.SetMode(config.Config.Mode)
	boot.engine = gin.Default()
	boot.regMiddlewares()
	boot.RegRoutes(boot.engine)
	boot.regRoutes()
	boot.engine.NoRoute(middleware.HandleNotFound)
	err := boot.engine.Run(config.Config.Address)
	db.Init(config.Config.Database, logger.NewLogger("mysql"))
	if err != nil {
		return
	}
}

func (boot *Bootstrap) regMiddlewares() {
	boot.engine.Use(middleware.Xid)
	for _, midFunc := range boot.HttpMiddlewares {
		boot.engine.Use(midFunc())
	}
	boot.engine.Use(gin.Recovery())
}

func (boot *Bootstrap) regRoutes() {
	boot.engine.GET("/metrics", middleware.PromHandler(promhttp.Handler()))
	boot.engine.GET("/consul", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}

func ginLogger() io.Writer {
	f, _ := os.Create("gin.log")
	return io.MultiWriter(f, os.Stdout)
}

func Response(ctx *gin.Context, result interface{}, code consts.ErrCode, err error) {
	response := gin.H{
		"code": code,
	}
	if result != nil {
		response["payload"] = result
	}
	code = consts.ErrCode(code)
	if err != nil {
		if config.Config.Mode != "release" {
			response["message"] = err.Error()
		} else {
			response["message"] = code.String()
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func Success(ctx *gin.Context, result interface{}) {
	Response(ctx, result, 0, nil)
}

func Fail(ctx *gin.Context, code consts.ErrCode, message string) {
	Response(ctx, nil, code, errors.New(message))
	ctx.Abort()
}
