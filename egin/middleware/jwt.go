package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/daodao97/egin/lib"
	"github.com/daodao97/egin/service/user"
	"github.com/daodao97/egin/utils"
	"github.com/daodao97/egin/utils/config"
)

func jwtAbort(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code":    http.StatusUnauthorized,
		"message": msg,
	})
	c.Abort()
}

func JWTMiddleware(u user.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, has := lib.Find(config.Config.Jwt.OpenApi, c.Request.URL.Path); has {
			c.Next()
			return
		}
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			jwtAbort(c, "Authorization Failed.")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			jwtAbort(c, "Authorization Failed.")
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			jwtAbort(c, "无效的Token "+err.Error())
			return
		}

		if time.Now().Unix() > claims.ExpiresAt {
			jwtAbort(c, "Token已过期")
			return
		}

		info, err := u.Info(claims.UserID)
		if err != nil  {
			jwtAbort(c, "用户状态异常" + err.Error())
			return
		}

		if info.Id != claims.UserID {
			jwtAbort(c, "无效的Token")
			return
		}

		c.Set("user", info)
		c.Next()
	}
}
