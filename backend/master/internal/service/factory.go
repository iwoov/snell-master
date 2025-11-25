package service

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/repository"
	"github.com/iwoov/snell-master/pkg/config"
)

// Services 聚合所有业务服务。
type Services struct {
	Admin        *AdminService
	User         *UserService
	Node         *NodeService
	Instance     *InstanceService
	Traffic      *TrafficService
	Subscribe    *SubscribeService
	Template     *TemplateService
	Log          *LogService
	Dashboard    *DashboardService
	SystemConfig *SystemConfigService
}

// ServiceDeps 注入依赖。
type ServiceDeps struct {
	Repositories *repository.Repositories
	Logger       *logrus.Logger
	Config       *config.Config
	DB           *gorm.DB
}

// NewServices 使用仓储和配置创建所有服务。
func NewServices(deps ServiceDeps) *Services {
	repos := deps.Repositories
	adminSvc := NewAdminService(repos.Admin, deps.Logger, deps.Config.JWT.Secret, deps.Config.JWT.ExpireHours)
	userSvc := NewUserService(repos.User, deps.Logger, deps.Config.JWT.Secret, deps.Config.JWT.ExpireHours)
	nodeSvc := NewNodeService(repos.Node, repos.Instance, deps.Logger)
	instanceSvc := NewInstanceService(repos.Instance, repos.User, repos.Node, deps.Logger)
	trafficSvc := NewTrafficService(repos.Traffic, repos.User, deps.Logger)
	subscribeSvc := NewSubscribeService(repos.Subscribe, repos.Template, repos.User, repos.Node, repos.Instance, deps.Logger)
	templateSvc := NewTemplateService(repos.Template, deps.Logger)
	logSvc := NewLogService(repos.Log, deps.Logger)
	dashboardSvc := NewDashboardService(deps.DB)
	systemConfigSvc := NewSystemConfigService(repos.SystemConfig, deps.Logger)

	return &Services{
		Admin:        adminSvc,
		User:         userSvc,
		Node:         nodeSvc,
		Instance:     instanceSvc,
		Traffic:      trafficSvc,
		Subscribe:    subscribeSvc,
		Template:     templateSvc,
		Log:          logSvc,
		Dashboard:    dashboardSvc,
		SystemConfig: systemConfigSvc,
	}
}
