package agent

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/middleware"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// Handler 提供 Agent 相关接口。
type Handler struct {
	nodeSvc     *service.NodeService
	instanceSvc *service.InstanceService
	trafficSvc  *service.TrafficService
}

// NewHandler 构造函数。
func NewHandler(nodeSvc *service.NodeService, instanceSvc *service.InstanceService, trafficSvc *service.TrafficService) *Handler {
	return &Handler{nodeSvc: nodeSvc, instanceSvc: instanceSvc, trafficSvc: trafficSvc}
}

// GetConfig 返回节点上的实例配置。
func (h *Handler) GetConfig(c *gin.Context) {
	node := middleware.GetNode(c)
	if node == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "node unavailable"})
		return
	}
	instances, err := h.instanceSvc.GetInstancesByNode(node.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	result := ConfigResponse{Instances: make([]InstanceConfig, 0, len(instances))}
	for _, inst := range instances {
		result.Instances = append(result.Instances, InstanceConfig{
			InstanceID: inst.ID,
			UserID:     inst.UserID,
			NodeID:     inst.NodeID,
			Port:       inst.Port,
			PSK:        inst.PSK,
			Version:    inst.Version,
			Obfs:       inst.Obfs,
		})
	}
	c.JSON(http.StatusOK, result)
}

// Heartbeat 上报心跳。
func (h *Handler) Heartbeat(c *gin.Context) {
	node := middleware.GetNode(c)
	if node == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "node unavailable"})
		return
	}
	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	if err := h.nodeSvc.UpdateHeartbeat(node.APIToken, req.CPUUsage, req.MemoryUsage, req.InstanceCount, req.Status, req.Version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ReportTraffic 批量上报流量。
func (h *Handler) ReportTraffic(c *gin.Context) {
	var req TrafficReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	now := time.Now()
	for _, record := range req.Records {
		if err := h.trafficSvc.RecordTraffic(record.UserID, record.InstanceID, record.NodeID, record.Upload, record.Download, now); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ReportInstanceStatus 更新实例状态。
func (h *Handler) ReportInstanceStatus(c *gin.Context) {
	var req StatusReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	for _, item := range req.Instances {
		if err := h.instanceSvc.UpdateInstanceStatus(item.InstanceID, item.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
