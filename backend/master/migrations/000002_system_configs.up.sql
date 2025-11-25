-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    description TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_system_configs_key ON system_configs(key);

-- 插入默认配置数据
INSERT INTO system_configs (key, value, description) VALUES
('snell_version', '5.0.1', 'Snell Server 版本号'),
('snell_base_url', 'https://dl.nssurge.com/snell', 'Snell Server 官方下载源'),
('snell_mirror_url', '', 'Snell Server 镜像下载源（可选，用于国内加速）'),
('agent_version', '1.0.0', 'Agent 版本号'),
('agent_download_url', 'https://github.com/iwoov/snell-master/releases/download/v{version}/snell-agent-{arch}', 'Agent 下载地址模板'),
('master_url', 'http://localhost:8080', 'Master 服务器地址（用于生成部署脚本）'),
('default_port_start', '10000', '默认端口范围起始'),
('default_port_end', '20000', '默认端口范围结束');
