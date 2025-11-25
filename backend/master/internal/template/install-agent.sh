#!/bin/bash
set -e

# ============================================
# Snell Master Agent 一键部署脚本
# 自动生成，请勿手动修改
# ============================================

# 配置变量（由 Master 自动填充）
MASTER_URL="{{.MasterURL}}"
API_TOKEN="{{.APIToken}}"
AGENT_VERSION="{{.AgentVersion}}"
NODE_NAME="{{.NodeName}}"
AGENT_DOWNLOAD_URL_TEMPLATE="{{.AgentDownloadURL}}"
AGENT_BINARY_URL_TEMPLATE="{{.AgentBinaryURL}}"

# 安装目录
INSTALL_DIR="/opt/snell-master/agent"
CONFIG_DIR="/etc/snell-master"
LOG_DIR="/var/log/snell-master"
DATA_DIR="/var/lib/snell-master"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}  Snell Master Agent 一键部署脚本${NC}"
echo -e "${GREEN}  节点: ${NODE_NAME}${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""

# 检查 root 权限
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}错误: 请使用 root 权限运行此脚本${NC}"
    exit 1
fi

# 检测系统架构
detect_arch() {
    local arch=$(uname -m)
    case $arch in
        x86_64)
            echo "amd64"
            ;;
        i686|i386)
            echo "i386"
            ;;
        aarch64|arm64)
            echo "aarch64"
            ;;
        armv7l)
            echo "armv7l"
            ;;
        *)
            echo -e "${RED}不支持的架构: $arch${NC}"
            exit 1
            ;;
    esac
}

# 检测操作系统
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo $ID
    else
        echo "unknown"
    fi
}

ARCH=$(detect_arch)
OS=$(detect_os)

echo -e "${GREEN}系统信息:${NC}"
echo "  操作系统: $OS"
echo "  架构: $ARCH"
echo ""

# 创建目录
echo -e "${YELLOW}[1/6] 创建目录...${NC}"
mkdir -p $INSTALL_DIR
mkdir -p $CONFIG_DIR
mkdir -p $LOG_DIR
mkdir -p $DATA_DIR/instances
echo -e "${GREEN}✓ 目录创建完成${NC}"
echo ""

# 下载 Agent
echo -e "${YELLOW}[2/6] 下载 Snell Master Agent...${NC}"
# 替换模板中的占位符
AGENT_URL=$(echo "$AGENT_DOWNLOAD_URL_TEMPLATE" | sed "s/{version}/$AGENT_VERSION/g" | sed "s/{arch}/$ARCH/g")
echo "下载地址: $AGENT_URL"

if command -v wget &> /dev/null; then
    wget -O $INSTALL_DIR/snell-agent $AGENT_URL
elif command -v curl &> /dev/null; then
    curl -L -o $INSTALL_DIR/snell-agent $AGENT_URL
else
    echo -e "${RED}错误: 需要安装 wget 或 curl${NC}"
    exit 1
fi

chmod +x $INSTALL_DIR/snell-agent
echo -e "${GREEN}✓ Agent 下载完成${NC}"
echo ""

# 生成配置文件
echo -e "${YELLOW}[3/6] 生成配置文件...${NC}"
cat > $CONFIG_DIR/agent.yaml << EOF
# Snell Master Agent 配置文件
agent:
  node_name: "${NODE_NAME}"
  master_url: "${MASTER_URL}"
  api_token: "${API_TOKEN}"

  # 实例管理
  instance_dir: "$DATA_DIR/instances"
  port_range_start: 10000
  port_range_end: 20000

  # Snell 二进制路径（Agent 会自动下载）
  snell_binary: "/usr/local/bin/snell-server"

  # 定时任务间隔（秒）
  heartbeat_interval: 30
  config_sync_interval: 60
  traffic_report_interval: 300

  # 日志
  log_level: "info"
  log_format: "json"
  log_file: "$LOG_DIR/agent.log"

monitor:
  enable_cpu: true
  enable_memory: true
  enable_traffic: true
EOF

echo -e "${GREEN}✓ 配置文件生成完成${NC}"
echo ""

# 安装 nftables（如果未安装）
echo -e "${YELLOW}[4/6] 检查 nftables...${NC}"
if ! command -v nft &> /dev/null; then
    echo "安装 nftables..."
    case $OS in
        ubuntu|debian)
            apt-get update && apt-get install -y nftables
            ;;
        centos|rhel|fedora)
            yum install -y nftables
            ;;
        *)
            echo -e "${YELLOW}警告: 请手动安装 nftables${NC}"
            ;;
    esac
fi
systemctl enable nftables 2>/dev/null || true
systemctl start nftables 2>/dev/null || true
echo -e "${GREEN}✓ nftables 准备完成${NC}"
echo ""

# 创建 systemd 服务
echo -e "${YELLOW}[5/6] 创建 systemd 服务...${NC}"
cat > /etc/systemd/system/snell-agent.service << EOF
[Unit]
Description=Snell Master Agent
After=network.target nftables.service

[Service]
Type=simple
User=root
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/snell-agent -c $CONFIG_DIR/agent.yaml
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal

# 安全选项
NoNewPrivileges=false
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable snell-agent
echo -e "${GREEN}✓ systemd 服务创建完成${NC}"
echo ""

# 启动服务
echo -e "${YELLOW}[6/6] 启动 Agent 服务...${NC}"
systemctl start snell-agent

# 等待服务启动
sleep 3

if systemctl is-active --quiet snell-agent; then
    echo -e "${GREEN}✓ Agent 服务启动成功${NC}"
else
    echo -e "${RED}✗ Agent 服务启动失败，请查看日志:${NC}"
    echo "  journalctl -u snell-agent -n 50"
    exit 1
fi

echo ""
echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}  部署完成！${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""
echo "查看服务状态:"
echo "  systemctl status snell-agent"
echo ""
echo "查看日志:"
echo "  journalctl -u snell-agent -f"
echo ""
echo "Agent 将自动从 Master 拉取配置并下载 Snell Server"
echo ""
