package public

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// AuthHandler 处理公开登录接口。
type AuthHandler struct {
	adminSvc *service.AdminService
	userSvc  *service.UserService
}

// NewAuthHandler 构造函数。
func NewAuthHandler(adminSvc *service.AdminService, userSvc *service.UserService) *AuthHandler {
	return &AuthHandler{adminSvc: adminSvc, userSvc: userSvc}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLogin 管理员登录。
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	token, admin, err := h.adminSvc.Login(req.Username, req.Password)
	if err != nil {
		common.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	common.Success(c, gin.H{"token": token, "admin": admin})
}

// UserLogin 用户登录。
func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	token, user, err := h.userSvc.Login(req.Username, req.Password)
	if err != nil {
		common.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}
	common.Success(c, gin.H{"token": token, "user": user})
}
