#!/bin/bash

set -e

# ==============================
# Terminal Colors
# ==============================
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
echo "┌──────────────────────────────────────────────┐"
echo "│     Backup Management System Installer       │"
echo "└──────────────────────────────────────────────┘"
echo -e "${NC}"

# ==============================
# Dependency Installer
# ==============================
install_dep() {
    local bin=$1
    local pkg=${2:-$1}

    if ! command -v "$bin" &>/dev/null; then
        echo -e "${YELLOW}Installing dependency: $pkg...${NC}"
        sudo apt update -qq
        sudo apt install -y "$pkg" -qq
        echo -e "${GREEN}✓ $pkg installed${NC}"
    else
        echo -e "${GREEN}✓ $bin already available${NC}"
    fi
}

# ==============================
# Network Utilities
# ==============================
get_ip() {
    hostname -I | awk '{print $1}'
}

find_available_port() {
    local ports=(5173 3000 3001 8080)
    for port in "${ports[@]}"; do
        if ! lsof -Pi :"$port" -sTCP:LISTEN -t >/dev/null 2>&1; then
            echo "$port"
            return
        fi
    done
    echo "5173"
}

# ==============================
# Dependency Check
# ==============================
echo -e "${BLUE}Verifying system requirements...${NC}"

install_dep "go" "golang-go"
install_dep "node" "nodejs"
install_dep "npm"
install_dep "mysql" "default-mysql-client"

# ==============================
# Rclone Installation
# ==============================
if ! command -v rclone &>/dev/null; then
    echo -e "${YELLOW}Installing rclone...${NC}"
    curl -fsSL https://rclone.org/install.sh | sudo bash > /dev/null
    echo -e "${GREEN}✓ rclone installed${NC}"
else
    echo -e "${GREEN}✓ rclone already available${NC}"
fi

# ==============================
# Database Configuration
# ==============================
echo ""
echo -e "${BLUE}Database Configuration${NC}"
echo ""

read -p "Database user [root]: " db_user
db_user=${db_user:-root}

read -sp "Database password: " db_pass
echo ""

read -p "Database host [127.0.0.1]: " db_host
db_host=${db_host:-127.0.0.1}

read -p "Database port [3306]: " db_port
db_port=${db_port:-3306}

read -p "Database name [gbackup_db]: " db_name
db_name=${db_name:-gbackup_db}

# ==============================
# Storage Configuration
# ==============================
echo ""
echo -e "${BLUE} Rclone Configuration${NC}"
echo ""

read -p "Rclone config file path [~/.config/rclone/rclone.conf]: " rclone_path
rclone_path=${rclone_path:-$HOME/.config/rclone/rclone.conf}

# ==============================
# Service Configuration
# ==============================
echo ""
echo -e "${BLUE}Service Configuration${NC}"
echo ""

APP_PORT=$(find_available_port)
jwt_secret=$(openssl rand -base64 32)

echo -e "${GREEN}✓ Application port selected: ${APP_PORT}${NC}"

# ==============================
# Environment Setup
# ==============================
echo ""
echo -e "${BLUE}Generating system configuration...${NC}"

cat > Backend/.env << EOF
DB_USER=$db_user
DB_PASS=$db_pass
DB_HOST=$db_host
DB_PORT=$db_port
DB_NAME=$db_name
JWT_SECRET_KEY=$jwt_secret
RCLONE_CONFIG=$rclone_path
SERVER_HOST=0.0.0.0
APP_PORT=$APP_PORT
EOF

echo -e "${GREEN}✓ Configuration created${NC}"

# ==============================
# Build Core Service
# ==============================
echo -e "${BLUE}Compiling application service...${NC}"
cd Backend
go mod download
go build -o app ./cmd/backend_app
cd ..
echo -e "${GREEN}✓ Application compiled${NC}"

# ==============================
# Install Web Interface
# ==============================
if [ -d "Frontend" ] && [ -f "Frontend/package.json" ]; then
    echo -e "${BLUE}Preparing web interface...${NC}"
    cd Frontend
    npm install --silent
    cd ..
    echo -e "${GREEN}✓ Web interface ready${NC}"
fi

# ==============================
# Completion Summary
# ==============================
SERVER_IP=$(get_ip)

echo ""
echo -e "${BLUE}┌──────────────────────────────────────────────┐"
echo -e "│          Installation Completed              │"
echo -e "└──────────────────────────────────────────────┘${NC}"
echo ""
echo -e "${GREEN}System successfully installed.${NC}"
echo ""
echo -e "${YELLOW}Access URL:${NC}"
echo -e "  http://${SERVER_IP}:${APP_PORT}"
echo ""
echo -e "${BLUE}To start the system:${NC}"
echo "  ./start.sh"
echo ""
