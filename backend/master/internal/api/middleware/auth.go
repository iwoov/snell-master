package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
	jwtutil "github.com/iwoov/snell-master/pkg/jwt"
)

const (
	contextUserIDKey   = "ctx_user_id"
	contextUsernameKey = "ctx_username"
	contextRoleKey     = "ctx_role"
	contextNodeIDKey   = "ctx_node_id"
	contextNodeKey     = "ctx_node"
)

// AuthMiddleware 解析 Authorization Bearer Token 并写入上下文。
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
			return
		}

		claims, err := jwtutil.ParseToken(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		c.Set(contextUserIDKey, claims.UserID)
		c.Set(contextUsernameKey, claims.Username)
		c.Set(contextRoleKey, claims.Role)
		c.Next()
	}
}

// AgentAuth 验证节点 API Token。
func AgentAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiToken := strings.TrimSpace(c.GetHeader("X-API-Token"))
		if apiToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing api token"})
			return
		}
		if db == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "database unavailable"})
			return
		}

		var node model.Node
		if err := db.Where("api_token = ?", apiToken).First(&node).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid api token"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "query node failed"})
			return
		}

		c.Set(contextNodeIDKey, node.ID)
		c.Set(contextNodeKey, &node)
		c.Next()
	}
}

// RequireAdmin 仅允许管理员访问。
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch GetUserRole(c) {
		case "admin", "super_admin":
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "admin role required"})
		}
	}
}

// RequireUser 要求普通用户身份。
func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserRole(c) != "user" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "user role required"})
			return
		}
		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID。
func GetUserID(c *gin.Context) uint {
	if v, ok := c.Get(contextUserIDKey); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

// GetUsername 返回用户名。
func GetUsername(c *gin.Context) string {
	if v, ok := c.Get(contextUsernameKey); ok {
		if name, ok := v.(string); ok {
			return name
		}
	}
	return ""
}

// GetUserRole 返回角色。
func GetUserRole(c *gin.Context) string {
	if v, ok := c.Get(contextRoleKey); ok {
		if role, ok := v.(string); ok {
			return role
		}
	}
	return ""
}

// GetNodeID 返回 Agent 节点 ID。
func GetNodeID(c *gin.Context) uint {
	if v, ok := c.Get(contextNodeIDKey); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

// GetNode 返回节点详情。
func GetNode(c *gin.Context) *model.Node {
	if v, ok := c.Get(contextNodeKey); ok {
		if node, ok := v.(*model.Node); ok {
			return node
		}
	}
	return nil
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}

	return authHeader
}
