package public

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// SubscribeHandler 提供订阅内容。
type SubscribeHandler struct {
	subscribeSvc *service.SubscribeService
}

// NewSubscribeHandler 构造函数。
func NewSubscribeHandler(subscribeSvc *service.SubscribeService) *SubscribeHandler {
	return &SubscribeHandler{subscribeSvc: subscribeSvc}
}

// GetSurgeSubscription 返回 Surge 配置。
func (h *SubscribeHandler) GetSurgeSubscription(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "token required"})
		return
	}
	content, err := h.subscribeSvc.GenerateSurgeConfig(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, content)
}
