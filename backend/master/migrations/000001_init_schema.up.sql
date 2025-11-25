-- Enable foreign keys for SQLite
PRAGMA foreign_keys = ON;

-- Admins table
CREATE TABLE IF NOT EXISTS admins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email TEXT,
    role INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email TEXT,
    traffic_limit BIGINT NOT NULL DEFAULT 107374182400,
    traffic_used_today BIGINT NOT NULL DEFAULT 0,
    traffic_used_month BIGINT NOT NULL DEFAULT 0,
    traffic_used_total BIGINT NOT NULL DEFAULT 0,
    reset_day INTEGER NOT NULL DEFAULT 1,
    status INTEGER NOT NULL DEFAULT 1,
    expire_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Nodes table
CREATE TABLE IF NOT EXISTS nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    api_token TEXT NOT NULL UNIQUE,
    endpoint TEXT NOT NULL,
    location TEXT,
    country_code TEXT,
    status TEXT NOT NULL DEFAULT 'offline',
    cpu_usage REAL DEFAULT 0,
    memory_usage REAL DEFAULT 0,
    disk_usage REAL DEFAULT 0,
    bandwidth_usage REAL DEFAULT 0,
    last_seen_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- User-Nodes relation
CREATE TABLE IF NOT EXISTS user_nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    connected_at DATETIME,
    UNIQUE(user_id, node_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- Snell instances
CREATE TABLE IF NOT EXISTS snell_instances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    port INTEGER NOT NULL,
    psk TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 4,
    obfs TEXT,
    config_path TEXT,
    service_name TEXT,
    status TEXT NOT NULL DEFAULT 'stopped',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(node_id, port),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- Traffic records
CREATE TABLE IF NOT EXISTS traffic_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    instance_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    upload_bytes BIGINT NOT NULL DEFAULT 0,
    download_bytes BIGINT NOT NULL DEFAULT 0,
    bytes_total BIGINT NOT NULL DEFAULT 0,
    record_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(instance_id) REFERENCES snell_instances(id) ON DELETE CASCADE,
    FOREIGN KEY(node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- Node heartbeats
CREATE TABLE IF NOT EXISTS node_heartbeats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER NOT NULL,
    status TEXT NOT NULL,
    message TEXT,
    cpu_usage REAL DEFAULT 0,
    memory_usage REAL DEFAULT 0,
    instance_count INTEGER DEFAULT 0,
    version TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- Templates
CREATE TABLE IF NOT EXISTS templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    content TEXT NOT NULL,
    is_default INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Subscribe tokens
CREATE TABLE IF NOT EXISTS subscribe_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    template_id INTEGER,
    token TEXT NOT NULL UNIQUE,
    expires_at DATETIME,
    last_access_at DATETIME,
    access_count INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(template_id) REFERENCES templates(id) ON DELETE SET NULL
);

-- Operation logs
CREATE TABLE IF NOT EXISTS operation_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    admin_id INTEGER,
    action TEXT NOT NULL,
    target_type TEXT,
    target_id INTEGER,
    details TEXT,
    ip_address TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(admin_id) REFERENCES admins(id) ON DELETE SET NULL
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes(status);
CREATE INDEX IF NOT EXISTS idx_user_nodes_user ON user_nodes(user_id);
CREATE INDEX IF NOT EXISTS idx_user_nodes_node ON user_nodes(node_id);
CREATE INDEX IF NOT EXISTS idx_instances_user ON snell_instances(user_id);
CREATE INDEX IF NOT EXISTS idx_instances_node ON snell_instances(node_id);
CREATE INDEX IF NOT EXISTS idx_traffic_user ON traffic_records(user_id);
CREATE INDEX IF NOT EXISTS idx_traffic_instance ON traffic_records(instance_id);
CREATE INDEX IF NOT EXISTS idx_traffic_node ON traffic_records(node_id);
CREATE INDEX IF NOT EXISTS idx_traffic_date ON traffic_records(record_date);
CREATE INDEX IF NOT EXISTS idx_heartbeats_node ON node_heartbeats(node_id);
CREATE INDEX IF NOT EXISTS idx_subscribe_tokens_token ON subscribe_tokens(token);
CREATE INDEX IF NOT EXISTS idx_operation_logs_admin ON operation_logs(admin_id);

-- Default admin (password: admin123)
INSERT INTO admins (username, password_hash, email, role)
VALUES (
    'admin',
    '$2a$10$KkXM9AfwI3Ct0KrQhg6O4eRTupp6JTDqKQRX4vISxwv3LxDjv4N7i',
    'admin@example.com',
    2
)
ON CONFLICT(username) DO NOTHING;

-- Default template
INSERT INTO templates (name, description, content, is_default)
VALUES (
    'default_surge',
    '默认 Surge 订阅模板',
    '[Proxy]\n{{node_list}}\n\n[Proxy Group]\nProxy = select, {{node_names}}\n\n[Rule]\nFINAL, Proxy',
    1
)
ON CONFLICT(name) DO NOTHING;
