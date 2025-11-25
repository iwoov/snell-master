package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// TemplateHandler 管理 Surge 模板。
type TemplateHandler struct {
	svc *service.TemplateService
}

// NewTemplateHandler 构造函数。
func NewTemplateHandler(svc *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{svc: svc}
}

// List 返回模板。
func (h *TemplateHandler) List(c *gin.Context) {
	items, err := h.svc.ListTemplates()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, items)
}

// Create 新建模板。
func (h *TemplateHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Content     string `json:"content" binding:"required"`
		Description string `json:"description"`
		IsDefault   bool   `json:"is_default"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	tpl, err := h.svc.CreateTemplate(req.Name, req.Content, req.Description, req.IsDefault)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Created(c, tpl)
}

// Update 编辑模板。
func (h *TemplateHandler) Update(c *gin.Context) {
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
	tpl, err := h.svc.UpdateTemplate(uint(id), updates)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, tpl)
}

// Delete 删除。
func (h *TemplateHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteTemplate(uint(id)); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"deleted": id})
}

// SetDefault 设置默认模板。
func (h *TemplateHandler) SetDefault(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.SetDefaultTemplate(uint(id)); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"default": id})
}
