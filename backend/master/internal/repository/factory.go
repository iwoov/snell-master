package repository

import "gorm.io/gorm"

// Repositories 聚合所有仓储依赖，便于注入。
type Repositories struct {
	Admin        AdminRepository
	User         UserRepository
	Node         NodeRepository
	Instance     InstanceRepository
	Traffic      TrafficRepository
	Subscribe    SubscribeRepository
	Template     TemplateRepository
	Log          LogRepository
	SystemConfig SystemConfigRepository
}

// NewRepositories 根据数据库实例创建所有仓储。
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Admin:        NewAdminRepository(db),
		User:         NewUserRepository(db),
		Node:         NewNodeRepository(db),
		Instance:     NewInstanceRepository(db),
		Traffic:      NewTrafficRepository(db),
		Subscribe:    NewSubscribeRepository(db),
		Template:     NewTemplateRepository(db),
		Log:          NewLogRepository(db),
		SystemConfig: NewSystemConfigRepository(db),
	}
}
