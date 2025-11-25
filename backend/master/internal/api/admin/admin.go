package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// AdminHandler 管理管理员账号。
type AdminHandler struct {
	svc *service.AdminService
}

// NewAdminHandler 构造函数。
func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// List 管理员列表。
func (h *AdminHandler) List(c *gin.Context) {
	page, pageSize := parsePagination(c)
	admins, total, err := h.svc.ListAdmins(page, pageSize)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, gin.H{"items": admins, "total": total, "page": page, "page_size": pageSize})
}

// Get 返回单个管理员。
func (h *AdminHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	admin, err := h.svc.GetAdminByID(uint(id))
	if err != nil {
		common.Fail(c, http.StatusNotFound, err.Error())
		return
	}
	common.Success(c, admin)
}

// Create 创建管理员。
func (h *AdminHandler) Create(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email"`
		Role     int    `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	admin, err := h.svc.CreateAdmin(req.Username, req.Password, req.Email, req.Role)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Created(c, admin)
}

// Update 更新管理员信息。
func (h *AdminHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Email string `json:"email"`
		Role  int    `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	admin, err := h.svc.UpdateAdmin(uint(id), req.Email, req.Role)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, admin)
}

// Delete 删除管理员。
func (h *AdminHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteAdmin(uint(id)); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"deleted": id})
}

// ChangePassword 修改管理员密码。
func (h *AdminHandler) ChangePassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.ChangePassword(uint(id), req.OldPassword, req.NewPassword); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"updated": id})
}

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return page, pageSize
}
