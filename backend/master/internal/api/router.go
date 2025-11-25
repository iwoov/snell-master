package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/api/handler"
	"github.com/iwoov/snell-master/backend/master/internal/api/middleware"
	"github.com/iwoov/snell-master/backend/master/internal/service"
	"github.com/iwoov/snell-master/pkg/config"
)

// SetupRouter 注册所有路由并返回 gin Engine。
func SetupRouter(cfg *config.Config, handlers *handler.Handlers, logSvc *service.LogService, db *gorm.DB) *gin.Engine {
	switch cfg.Server.Mode {
	case gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(middleware.Logger(cfg.Log.IgnorePaths), middleware.CORS(cfg.CORS), gin.Recovery())

	r.GET("/healthz", handlers.Health.Health)

	publicGroup := r.Group("/api")
	{
		publicGroup.POST("/auth/admin/login", handlers.Auth.AdminLogin)
		publicGroup.POST("/auth/user/login", handlers.Auth.UserLogin)
		publicGroup.GET("/subscribe/:token", handlers.PublicSubscribe.GetSurgeSubscription)
		publicGroup.GET("/health", handlers.Health.Health)
		publicGroup.GET("/ping", handlers.Health.Ping)
	}

	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	adminGroup.Use(middleware.RequireAdmin())
	adminGroup.Use(middleware.OperationLogger(logSvc))
	{
		admins := adminGroup.Group("/admins")
		admins.GET("", handlers.Admin.List)
		admins.POST("", handlers.Admin.Create)
		admins.GET("/:id", handlers.Admin.Get)
		admins.PUT("/:id", handlers.Admin.Update)
		admins.DELETE("/:id", handlers.Admin.Delete)
		admins.POST("/:id/password", handlers.Admin.ChangePassword)

		users := adminGroup.Group("/users")
		users.GET("", handlers.AdminUser.List)
		users.POST("", handlers.AdminUser.Create)
		users.GET("/:id", handlers.AdminUser.Get)
		users.PUT("/:id", handlers.AdminUser.Update)
		users.DELETE("/:id", handlers.AdminUser.Delete)
		users.POST("/:id/reset", handlers.AdminUser.ResetTraffic)
		users.POST("/:id/status", handlers.AdminUser.UpdateStatus)
		users.POST("/:id/nodes", handlers.AdminUser.AssignNodes)

		nodes := adminGroup.Group("/nodes")
		nodes.GET("", handlers.Node.List)
		nodes.POST("", handlers.Node.Create)
		nodes.GET("/:id", handlers.Node.Get)
		nodes.PUT("/:id", handlers.Node.Update)
		nodes.DELETE("/:id", handlers.Node.Delete)
		nodes.POST("/:id/token", handlers.Node.RegenerateToken)
		nodes.GET("/:id/install-script", handlers.Node.GetInstallScript)

		instances := adminGroup.Group("/instances")
		instances.GET("", handlers.Instance.List)
		instances.POST("", handlers.Instance.Create)
		instances.GET("/:id", handlers.Instance.Get)
		instances.DELETE("/:id", handlers.Instance.Delete)
		instances.PUT("/:id/status", handlers.Instance.UpdateStatus)
		instances.POST("/:id/restart", handlers.Instance.Restart)

		traffic := adminGroup.Group("/traffic")
		traffic.GET("/summary", handlers.Traffic.Summary)
		traffic.GET("/users", handlers.Traffic.UserRanking)
		traffic.GET("/trend", handlers.Traffic.LineTrend)
		traffic.GET("/nodes/:id", handlers.Traffic.NodeTraffic)

		subs := adminGroup.Group("/subscriptions")
		subs.GET("", handlers.Subscribe.List)
		subs.POST("", handlers.Subscribe.Create)
		subs.DELETE("/:id", handlers.Subscribe.Delete)
		subs.POST("/:id/regenerate", handlers.Subscribe.Regenerate)

		templates := adminGroup.Group("/templates")
		templates.GET("", handlers.Template.List)
		templates.POST("", handlers.Template.Create)
		templates.PUT("/:id", handlers.Template.Update)
		templates.DELETE("/:id", handlers.Template.Delete)
		templates.POST("/:id/default", handlers.Template.SetDefault)

		// 系统配置管理路由
		sysConfigs := adminGroup.Group("/system-configs")
		sysConfigs.GET("", handlers.SystemConfig.GetSystemConfigs)
		sysConfigs.PUT("/:key", handlers.SystemConfig.UpdateSystemConfig)
		sysConfigs.PUT("", handlers.SystemConfig.BatchUpdateConfigs)

		adminGroup.GET("/logs", handlers.Log.List)
		adminGroup.GET("/dashboard/stats", handlers.Dashboard.Stats)
	}

	userGroup := r.Group("/api/user")
	userGroup.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	userGroup.Use(middleware.RequireUser())
	{
		userGroup.GET("/profile", handlers.UserProfile.GetProfile)
		userGroup.PUT("/profile", handlers.UserProfile.UpdateProfile)
		userGroup.POST("/password", handlers.UserProfile.ChangePassword)
		userGroup.GET("/instances", handlers.UserInstance.ListMyInstances)
		userGroup.GET("/instances/:id", handlers.UserInstance.GetMyInstance)
		userGroup.GET("/traffic", handlers.UserTraffic.GetMyTraffic)
		userGroup.GET("/subscriptions", handlers.UserSubscribe.GetMySubscription)
		userGroup.POST("/subscriptions/regenerate", handlers.UserSubscribe.RegenerateToken)
	}

	agentGroup := r.Group("/api/agent")
	agentGroup.Use(middleware.AgentAuth(db))
	{
		agentGroup.GET("/config", handlers.Agent.GetConfig)
		agentGroup.POST("/heartbeat", handlers.Agent.Heartbeat)
		agentGroup.POST("/traffic", handlers.Agent.ReportTraffic)
		agentGroup.POST("/status", handlers.Agent.ReportInstanceStatus)
		agentGroup.GET("/snell-config", handlers.AgentSnell.GetSnellConfig)
	}

	return r
}
