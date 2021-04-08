package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/daodao97/egin/lib"
	"github.com/daodao97/egin/utils/logger"
)

// 日志记录到文件
func HttpLog() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 请求路由
		reqUri := c.Request.RequestURI
		if _, ok := lib.Find([]string{"/metrics"}, reqUri); ok {
			return
		}
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		reqId, _ := c.Get("req_id")
		// 日志格式
		logger.NewLogger("http").Info("http end", map[string]interface{}{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
			"req_id":       reqId,
		})
	}
}
