package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/daodao97/egin/egin/lib"
	"github.com/daodao97/egin/egin/utils"
	"github.com/daodao97/egin/egin/utils/config"
)

// cors
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.Config.Auth.Cors
		if !conf.Enable {
			return
		}
		allowOrigin := conf.AllowOrigins
		if len(allowOrigin) == 0 {
			allowOrigin = []string{"*"}
		}
		allCredentials := "true"
		if !conf.AllowCredentials {
			allCredentials = "false"
		}
		allowHeaders := []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
		allowMethods := []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"}
		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(allowOrigin, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", allCredentials)
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeaders, ","))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(allowMethods, ","))
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ip白名单
func IPAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.Config.Auth.IpAuth
		if !conf.Enable {
			return
		}
		clientIp := c.ClientIP()
		_, hasIt := lib.Find(conf.AllowedIpList, clientIp)
		if !hasIt {
			c.String(http.StatusUnauthorized, "%s, not in ipList", clientIp)
			c.Abort()
		}
		c.Next()
	}
}

// 基于 AK,SK 的签名验证
func MacAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.Config.Auth.AKSK
		if !conf.Enable {
			return
		}

		req := c.Request
		auth := req.Header.Get("Authorization")
		if auth == "" {
			c.String(http.StatusUnauthorized, "need auth")
			c.Abort()
			return
		}
		info := strings.Split(auth, ":")
		if len(info) != 2 {
			c.String(http.StatusUnauthorized, "auth error")
			c.Abort()
			return
		}
		accessKey := info[0]
		secretKey, ok := conf.Allowed[accessKey]
		if !ok {
			c.String(http.StatusUnauthorized, "auth user not found")
			c.Abort()
			return
		}
		mac := utils.Mac{AccessKey: accessKey, SecretKey: []byte(secretKey)}
		token, _ := mac.SignRequest(req)
		fmt.Println(token, auth)
		if token != auth {
			c.String(http.StatusUnauthorized, "Auth failed")
			c.Abort()
			return
		}
		c.Next()
	}
}
