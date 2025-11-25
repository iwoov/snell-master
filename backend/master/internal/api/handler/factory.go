package handler

import (
	"time"

	"gorm.io/gorm"

	adminapi "github.com/iwoov/snell-master/backend/master/internal/api/admin"
	agentapi "github.com/iwoov/snell-master/backend/master/internal/api/agent"
	publicapi "github.com/iwoov/snell-master/backend/master/internal/api/public"
	userapi "github.com/iwoov/snell-master/backend/master/internal/api/user"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// Handlers 聚合所有 HTTP 处理器。
type Handlers struct {
	Auth            *publicapi.AuthHandler
	Health          *publicapi.HealthHandler
	Admin           *adminapi.AdminHandler
	AdminUser       *adminapi.UserHandler
	Node            *adminapi.NodeHandler
	Instance        *adminapi.InstanceHandler
	Traffic         *adminapi.TrafficHandler
	Subscribe       *adminapi.SubscribeHandler
	Template        *adminapi.TemplateHandler
	Log             *adminapi.LogHandler
	Dashboard       *adminapi.DashboardHandler
	SystemConfig    *adminapi.SystemConfigHandler
	UserProfile     *userapi.ProfileHandler
	UserInstance    *userapi.InstanceHandler
	UserTraffic     *userapi.TrafficHandler
	UserSubscribe   *userapi.SubscribeHandler
	Agent           *agentapi.Handler
	AgentSnell      *agentapi.SnellHandler
	PublicSubscribe *publicapi.SubscribeHandler
}

// NewHandlers 初始化所有 Handler。
func NewHandlers(services *service.Services, db *gorm.DB, startTime time.Time) *Handlers {
	return &Handlers{
		Auth:            publicapi.NewAuthHandler(services.Admin, services.User),
		Health:          publicapi.NewHealthHandler(db, startTime),
		Admin:           adminapi.NewAdminHandler(services.Admin),
		AdminUser:       adminapi.NewUserHandler(services.User),
		Node:            adminapi.NewNodeHandler(services.Node, services.SystemConfig),
		Instance:        adminapi.NewInstanceHandler(services.Instance),
		Traffic:         adminapi.NewTrafficHandler(services.Traffic),
		Subscribe:       adminapi.NewSubscribeHandler(services.Subscribe),
		Template:        adminapi.NewTemplateHandler(services.Template),
		Log:             adminapi.NewLogHandler(services.Log),
		Dashboard:       adminapi.NewDashboardHandler(services.Dashboard),
		SystemConfig:    adminapi.NewSystemConfigHandler(services.SystemConfig),
		UserProfile:     userapi.NewProfileHandler(services.User),
		UserInstance:    userapi.NewInstanceHandler(services.Instance),
		UserTraffic:     userapi.NewTrafficHandler(services.Traffic),
		UserSubscribe:   userapi.NewSubscribeHandler(services.Subscribe),
		Agent:           agentapi.NewHandler(services.Node, services.Instance, services.Traffic),
		AgentSnell:      agentapi.NewSnellHandler(services.SystemConfig),
		PublicSubscribe: publicapi.NewSubscribeHandler(services.Subscribe),
	}
}
