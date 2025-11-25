package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// OperationLogger 记录管理员修改操作。
func OperationLogger(logSvc *service.LogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if logSvc == nil {
			return
		}
		if c.Request.Method == http.MethodGet {
			return
		}
		if status := c.Writer.Status(); status >= 400 {
			return
		}
		adminID := GetUserID(c)
		if adminID == 0 {
			return
		}
		route := c.FullPath()
		method := c.Request.Method
		if err := logSvc.LogOperation(adminID, method+" "+route, "", nil, "", c.ClientIP()); err != nil {
			// ignore logging errors
		}
	}
}
