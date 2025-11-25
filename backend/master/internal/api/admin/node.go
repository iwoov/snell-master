package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/service"
	"github.com/iwoov/snell-master/backend/master/internal/template"
)

// NodeHandler 管理节点。
type NodeHandler struct {
	svc           *service.NodeService
	configService *service.SystemConfigService
}

// NewNodeHandler 构造函数。
func NewNodeHandler(svc *service.NodeService, configService *service.SystemConfigService) *NodeHandler {
	return &NodeHandler{
		svc:           svc,
		configService: configService,
	}
}

// List 返回所有节点。
func (h *NodeHandler) List(c *gin.Context) {
	nodes, err := h.svc.GetNodeList()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, nodes)
}

// Create 创建节点。
func (h *NodeHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Endpoint    string `json:"endpoint" binding:"required"`
		Location    string `json:"location"`
		CountryCode string `json:"country_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	node, err := h.svc.RegisterNode(req.Name, req.Endpoint, req.Location, req.CountryCode)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Created(c, node)
}

// Get 返回详情。
func (h *NodeHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	node, err := h.svc.GetNodeByID(uint(id))
	if err != nil {
		common.Fail(c, http.StatusNotFound, err.Error())
		return
	}
	common.Success(c, node)
}

// Update 修改信息。
func (h *NodeHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	node, err := h.svc.UpdateNode(uint(id), updates)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, node)
}

// Delete 删除节点。
func (h *NodeHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteNode(uint(id)); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"deleted": id})
}

// RegenerateToken 重新生成 API Token。
func (h *NodeHandler) RegenerateToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	token, err := h.svc.RegenerateToken(uint(id))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"token": token})
}

// GetInstallScript 生成并下载节点部署脚本
// GET /api/admin/nodes/:id/install-script
func (h *NodeHandler) GetInstallScript(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	// 获取节点信息
	node, err := h.svc.GetNodeByID(uint(id))
	if err != nil {
		common.Fail(c, http.StatusNotFound, "node not found")
		return
	}

	// 获取系统配置
	systemConfig, err := h.configService.GetSystemConfig()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "failed to get system config")
		return
	}

	// 构造脚本数据
	scriptData := template.InstallScriptData{
		MasterURL:        systemConfig.MasterURL,
		APIToken:         node.APIToken,
		AgentVersion:     systemConfig.AgentVersion,
		NodeName:         node.Name,
		AgentDownloadURL: systemConfig.AgentDownloadURL,
		AgentBinaryURL:   systemConfig.AgentDownloadURL, // 使用相同的URL模板
	}

	// 生成脚本
	script, err := template.GenerateInstallScript(scriptData)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "failed to generate script")
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/x-shellscript; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=install-agent-%s.sh", node.Name))

	// 返回脚本内容
	c.String(http.StatusOK, script)
}
