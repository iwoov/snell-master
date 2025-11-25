package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// SystemConfigHandler 系统配置处理器
type SystemConfigHandler struct {
	configService *service.SystemConfigService
}

// NewSystemConfigHandler 创建系统配置处理器实例
func NewSystemConfigHandler(configService *service.SystemConfigService) *SystemConfigHandler {
	return &SystemConfigHandler{
		configService: configService,
	}
}

// GetSystemConfigs 获取所有系统配置
// GET /api/admin/system-configs
func (h *SystemConfigHandler) GetSystemConfigs(c *gin.Context) {
	configs, err := h.configService.GetAllConfigs()
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Failed to get system configs",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    configs,
	})
}

// UpdateSystemConfig 更新单个系统配置
// PUT /api/admin/system-configs/:key
func (h *SystemConfigHandler) UpdateSystemConfig(c *gin.Context) {
	key := c.Param("key")

	var req struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.configService.UpdateSystemConfig(key, req.Value); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Failed to update config",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "Config updated successfully",
	})
}

// BatchUpdateConfigs 批量更新系统配置
// PUT /api/admin/system-configs
func (h *SystemConfigHandler) BatchUpdateConfigs(c *gin.Context) {
	var req map[string]string

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid request",
		})
		return
	}

	if err := h.configService.BatchUpdateConfigs(req); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Failed to batch update configs",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "Configs updated successfully",
	})
}
