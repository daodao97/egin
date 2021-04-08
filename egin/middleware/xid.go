package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func Xid(c *gin.Context) {
	c.Set("req_id", xid.New().String())
}
