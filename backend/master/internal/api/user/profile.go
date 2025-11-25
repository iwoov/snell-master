package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/api/middleware"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// ProfileHandler 处理用户资料接口。
type ProfileHandler struct {
	svc *service.UserService
}

// NewProfileHandler 构造函数。
func NewProfileHandler(svc *service.UserService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

// GetProfile 返回个人资料。
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.svc.GetUserByID(userID)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, user)
}

// UpdateProfile 更新资料。
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := h.svc.UpdateUser(userID, updates)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, user)
}

// ChangePassword 修改密码。
func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"updated": true})
}
