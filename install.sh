#!/bin/bash

set -e

# Colors for clean output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'  # No Color

echo -e "${BLUE}"
echo "┌──────────────────────────────────────────────┐"
echo "│         Backup Server Installation          │"
echo "└──────────────────────────────────────────────┘"
echo -e "${NC}"

# Function to check and install dependency
install_dep() {
    local dep=$1
    local pkg=${2:-$1}
    
    if ! command -v $dep &>/dev/null; then
        echo -e "${YELLOW}Installing $dep...${NC}"
        sudo apt update > /dev/null 2>&1
        sudo apt install -y $pkg > /dev/null 2>&1
        echo -e "${GREEN}✓ $dep installed${NC}"
    else
        echo -e "${GREEN}✓ $dep already installed${NC}"
    fi
}

# Get server IP
get_ip() {
    ip=$(hostname -I | awk '{print $1}')
    [ -z "$ip" ] && ip=$(ip addr show | grep -oP 'inet \K[\d.]+' | grep -v 127.0.0.1 | head -n1)
    echo $ip
}

echo -e "${BLUE}Checking dependencies...${NC}"
install_dep "go" "golang-go"
install_dep "node" "nodejs"
install_dep "npm"
install_dep "mysql" "mysql-client"

# Install rclone
if ! command -v rclone &>/dev/null; then
    echo -e "${YELLOW}Installing rclone...${NC}"
    curl -s https://rclone.org/install.sh | sudo bash > /dev/null 2>&1
    echo -e "${GREEN}✓ rclone installed${NC}"
else
    echo -e "${GREEN}✓ rclone already installed${NC}"
fi

echo ""
echo -e "${BLUE}Database Configuration${NC}"
echo ""

read -p "MySQL Username: " db_user
db_user=${db_user:-gbackup}

read -sp "MySQL Password: " db_pass
db_pass=${db_pass:-gbackup}
echo ""

read -p "Host: " db_host
db_host=${db_host:-localhost}

read -p "Port: " db_port
db_port=${db_port:-3306}

read -p "Database: " db_name
db_name=${db_name:-gbackup_db}

echo ""
echo -e "${BLUE}Rclone Configuration${NC}"
echo ""
read -p "Rclone config path [~/.config/rclone/rclone.conf]: " rclone_path
rclone_path=${rclone_path:-$HOME/.config/rclone/rclone.conf}

echo ""
echo -e "${BLUE}Network Configuration${NC}"
echo ""
read -p "Backend port [8080]: " backend_port
backend_port=${backend_port:-8080}

read -p "Frontend port [3000]: " frontend_port
frontend_port=${frontend_port:-3000}

# Generate JWT
jwt_secret=$(openssl rand -base64 32)

echo ""
echo -e "${BLUE}Creating configuration...${NC}"

# Create .env
cat > Backend/.env << EOF
DB_USER=$db_user
DB_PASS=$db_pass
DB_HOST=$db_host
DB_PORT=$db_port
DB_NAME=$db_name
JWT_SECRET_KEY=$jwt_secret
RCLONE_CONFIG=$rclone_path
SERVER_HOST=0.0.0.0
BACKEND_PORT=$backend_port
FRONTEND_PORT=$frontend_port
EOF

echo -e "${GREEN}✓ Configuration file created${NC}"

# Build backend
echo -e "${BLUE}Building backend...${NC}"
cd Backend
go mod download
go build -o cmd/backend_app ./cmd
cd ..
echo -e "${GREEN}✓ Backend compiled${NC}"

# Install frontend
if [ -d "Frontend" ] && [ -f "Frontend/package.json" ]; then
    echo -e "${BLUE}Installing frontend...${NC}"
    cd Frontend
    npm install --silent
    cd ..
    echo -e "${GREEN}✓ Frontend dependencies installed${NC}"
fi

echo ""
echo -e "${BLUE}┌──────────────────────────────────────────────┐"
echo -e "│           Installation Complete              │"
echo -e "└──────────────────────────────────────────────┘${NC}"
echo ""
echo -e "${GREEN}Server is ready to run.${NC}"
echo ""
echo "To start the server:"
echo -e "  ${BLUE}./start.sh${NC}"
echo ""