package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/pkg/config"
)

// CORS 根据配置允许跨域请求。
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	allowAll := len(cfg.AllowOrigins) == 0
	allowed := make(map[string]struct{})
	for _, origin := range cfg.AllowOrigins {
		if origin == "*" {
			allowAll = true
			break
		}
		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowedOrigin := ""

		switch {
		case allowAll && origin != "":
			allowedOrigin = origin
		case allowAll:
			allowedOrigin = "*"
		case origin != "":
			if _, ok := allowed[origin]; ok {
				allowedOrigin = origin
			}
		}

		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
