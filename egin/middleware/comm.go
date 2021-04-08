package middleware

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var proxyMap = map[string]string{
	"order": "http://127.0.0.1:8080",
}

// 路由转发
func HandleNotFound(c *gin.Context) {
	urlParts := strings.Split(c.Request.URL.Path, "/")
	target, ok := proxyMap[urlParts[1]]
	if !ok {
		return
	}
	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = "/" + strings.Join(urlParts[2:], "/")
		// 鉴权信息的准换
		// jwt -> aksk
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
