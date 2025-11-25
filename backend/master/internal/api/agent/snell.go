package agent

import (
	"github.com/gin-gonic/gin"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// SnellHandler Snell 配置处理器
type SnellHandler struct {
	configService *service.SystemConfigService
}

// NewSnellHandler 创建 Snell 配置处理器实例
func NewSnellHandler(configService *service.SystemConfigService) *SnellHandler {
	return &SnellHandler{
		configService: configService,
	}
}

// GetSnellConfig 获取 Snell Server 下载配置
// GET /api/agent/snell-config
func (h *SnellHandler) GetSnellConfig(c *gin.Context) {
	config, err := h.configService.GetSnellConfig()
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Failed to get Snell config",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}
