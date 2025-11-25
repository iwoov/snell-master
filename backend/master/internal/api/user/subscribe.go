package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/api/middleware"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// SubscribeHandler 用户订阅接口。
type SubscribeHandler struct {
	svc *service.SubscribeService
}

// NewSubscribeHandler 构造函数。
func NewSubscribeHandler(svc *service.SubscribeService) *SubscribeHandler {
	return &SubscribeHandler{svc: svc}
}

// GetMySubscription 返回当前订阅。
func (h *SubscribeHandler) GetMySubscription(c *gin.Context) {
	userID := middleware.GetUserID(c)
	token, err := h.svc.GetTokenByUser(userID)
	if err != nil {
		common.Fail(c, http.StatusNotFound, err.Error())
		return
	}
	common.Success(c, token)
}

// RegenerateToken 重新生成订阅令牌。
func (h *SubscribeHandler) RegenerateToken(c *gin.Context) {
	userID := middleware.GetUserID(c)
	token, err := h.svc.RegenerateToken(userID)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"token": token})
}
