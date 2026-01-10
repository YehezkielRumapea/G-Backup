#!/bin/bash

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}"
echo "┌──────────────────────────────────────────────┐"
echo "│        System Installation Wizard           │"
echo "└──────────────────────────────────────────────┘"
echo -e "${NC}"

# ============================================
# Dependency Check & Install
# ============================================
echo -e "${BLUE}System Dependencies Check${NC}"

install_dep() {
    if ! command -v $1 &>/dev/null; then
        echo -e "${YELLOW}Installing $1...${NC}"
        sudo apt update -qq
        sudo apt install -y $2 -qq
        echo -e "${GREEN}✓ $1 installed${NC}"
    else
        echo -e "${GREEN}✓ $1 available${NC}"
    fi
}

install_dep "go" "golang-go"
install_dep "node" "nodejs"
install_dep "npm"
install_dep "mysql" "default-mysql-client"

# Rclone
if ! command -v rclone &>/dev/null; then
    echo -e "${YELLOW}Installing rclone...${NC}"
    curl -fsSL https://rclone.org/install.sh | sudo bash > /dev/null
    echo -e "${GREEN}✓ rclone installed${NC}"
else
    echo -e "${GREEN}✓ rclone available${NC}"
fi

# ============================================
# MySQL Setup
# ============================================
echo -e "\n${BLUE}Database Server Setup${NC}"

check_mysql() {
    # Check if MySQL is running
    if systemctl is-active --quiet mysql 2>/dev/null || \
       systemctl is-active --quiet mysqld 2>/dev/null; then
        echo -e "${GREEN}✓ MySQL server running${NC}"
        return 0
    fi
    
    # Check if installed but not running
    if dpkg -l | grep -q mysql-server || dpkg -l | grep -q mariadb-server; then
        echo -e "${YELLOW}Starting MySQL service...${NC}"
        sudo systemctl start mysql 2>/dev/null || sudo systemctl start mysqld
        sleep 3
        return 0
    fi
    
    # MySQL not installed
    return 1
}

mysql_newly_installed=0
if check_mysql; then
    echo -e "${GREEN}✓ MySQL ready${NC}"
else
    # Install MySQL with user-defined credentials
    echo -e "${YELLOW}MySQL not found. Installing MySQL server...${NC}"
    echo ""
    echo -e "${BLUE}Configure MySQL Administrator Account${NC}"
    echo -e "${YELLOW}Set username and password for MySQL${NC}"
    
    # Get username from user
    read -p "MySQL username [root]: " MYSQL_NEW_USER
    MYSQL_NEW_USER=${MYSQL_NEW_USER:-root}
    
    # Get password from user
    while true; do
        read -sp "MySQL password for '$MYSQL_NEW_USER': " MYSQL_NEW_PASS
        echo ""
        
        if [ -z "$MYSQL_NEW_PASS" ]; then
            echo -e "${RED}Password cannot be empty${NC}"
            continue
        fi
        
        if [ ${#MYSQL_NEW_PASS} -lt 6 ]; then
            echo -e "${RED}Password must be at least 6 characters${NC}"
            continue
        fi
        
        read -sp "Confirm password: " MYSQL_NEW_PASS_CONFIRM
        echo ""
        
        if [ "$MYSQL_NEW_PASS" != "$MYSQL_NEW_PASS_CONFIRM" ]; then
            echo -e "${RED}Passwords do not match${NC}"
        else
            break
        fi
    done
    
    # Set debconf for non-interactive install
    echo "mysql-server mysql-server/root_password password $MYSQL_NEW_PASS" | sudo debconf-set-selections
    echo "mysql-server mysql-server/root_password_again password $MYSQL_NEW_PASS" | sudo debconf-set-selections
    
    echo ""
    echo -e "${YELLOW}Installing MySQL server...${NC}"
    sudo apt update -qq
    sudo DEBIAN_FRONTEND=noninteractive apt install -y mysql-server -qq
    
    sudo systemctl start mysql
    sudo systemctl enable mysql
    
    sleep 5
    echo -e "${GREEN}✓ MySQL server installed${NC}"
    
    # If user chose non-root username, create that user
    if [ "$MYSQL_NEW_USER" != "root" ]; then
        echo -e "${YELLOW}Creating user '$MYSQL_NEW_USER'...${NC}"
        mysql -u root -p"$MYSQL_NEW_PASS" -e "CREATE USER IF NOT EXISTS '$MYSQL_NEW_USER'@'localhost' IDENTIFIED BY '$MYSQL_NEW_PASS';" 2>/dev/null
        mysql -u root -p"$MYSQL_NEW_PASS" -e "GRANT ALL PRIVILEGES ON *.* TO '$MYSQL_NEW_USER'@'localhost' WITH GRANT OPTION;" 2>/dev/null
        mysql -u root -p"$MYSQL_NEW_PASS" -e "FLUSH PRIVILEGES;" 2>/dev/null
        echo -e "${GREEN}✓ User '$MYSQL_NEW_USER' created${NC}"
    fi
    
    DB_USER="$MYSQL_NEW_USER"
    DB_PASS="$MYSQL_NEW_PASS"
    mysql_newly_installed=1
fi

# ============================================
# MySQL Authentication
# ============================================
echo -e "\n${BLUE}MySQL User Credentials${NC}"

if [ $mysql_newly_installed -eq 1 ]; then
    echo -e "${GREEN}Using newly created account: $DB_USER${NC}"
else
    echo -e "${YELLOW}Enter your existing MySQL credentials${NC}"
    # Ask for existing MySQL credentials
    attempts=0
    max_attempts=3
    
    while [ $attempts -lt $max_attempts ]; do
        echo ""
        read -p "MySQL username [root]: " DB_USER
        DB_USER=${DB_USER:-root}
        
        read -sp "MySQL password: " DB_PASS
        echo ""
        
        if mysql -u "$DB_USER" -p"$DB_PASS" -e "SELECT 1;" 2>/dev/null; then
            echo -e "${GREEN}✓ Authentication successful${NC}"
            break
        else
            attempts=$((attempts + 1))
            remaining=$((max_attempts - attempts))
            if [ $remaining -gt 0 ]; then
                echo -e "${RED}✗ Invalid credentials ($remaining attempts left)${NC}"
            else
                echo -e "${RED}✗ Maximum attempts reached${NC}"
                exit 1
            fi
        fi
    done
fi

echo -e "${GREEN}✓ Will use MySQL user: $DB_USER${NC}"

# ============================================
# Database Configuration
# ============================================
echo -e "\n${BLUE}Database Configuration${NC}"

read -p "Database name [backup_system]: " DB_NAME
DB_NAME=${DB_NAME:-backup_system}

read -p "Database host [localhost]: " DB_HOST
DB_HOST=${DB_HOST:-localhost}

read -p "Database port [3306]: " DB_PORT
DB_PORT=${DB_PORT:-3306}

# ============================================
# Create Database
# ============================================
echo -e "\n${BLUE}Setting up database...${NC}"

# Create database if not exists
mysql -u "$DB_USER" -p"$DB_PASS" -e "CREATE DATABASE IF NOT EXISTS \`$DB_NAME\`;" 2>/dev/null
echo -e "${GREEN}✓ Database '$DB_NAME' created${NC}"

# Test connection
if mysql -u "$DB_USER" -p"$DB_PASS" -D "$DB_NAME" -e "SELECT 1;" 2>/dev/null; then
    echo -e "${GREEN}✓ Database connection test passed${NC}"
else
    echo -e "${RED}✗ Database connection failed${NC}"
    exit 1
fi

# ============================================
# Rclone Configuration
# ============================================
echo -e "\n${BLUE}Storage Configuration${NC}"

read -p "Rclone config path [~/.config/rclone/rclone.conf]: " RCLONE_PATH
RCLONE_PATH=${RCLONE_PATH:-$HOME/.config/rclone/rclone.conf}

# ============================================
# Service Port Configuration
# ============================================
echo -e "\n${BLUE}Service Configuration${NC}"

# Find available port
APP_PORT=""
for port in 8080 3000 3001 5173; do
    if ! lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        APP_PORT=$port
        break
    fi
done

if [ -z "$APP_PORT" ]; then
    echo -e "${YELLOW}All default ports are in use${NC}"
    read -p "Enter custom port: " APP_PORT
fi

# Generate JWT secret
JWT_SECRET=$(openssl rand -base64 32)

echo -e "${GREEN}✓ Service port: $APP_PORT${NC}"
echo -e "${GREEN}✓ Security token generated${NC}"

# ============================================
# Create Environment File
# ============================================
echo -e "\n${BLUE}Creating configuration file...${NC}"

cat > Backend/.env << EOF
# Database Configuration
DB_USER=$DB_USER
DB_PASS=$DB_PASS
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_NAME=$DB_NAME

# Security
JWT_SECRET_KEY=$JWT_SECRET

# Storage
RCLONE_CONFIG=$RCLONE_PATH

# Network
SERVER_HOST=0.0.0.0
APP_PORT=$APP_PORT
EOF

echo -e "${GREEN}✓ Configuration created${NC}"

# ============================================
# Build Backend Application
# ============================================
echo -e "\n${BLUE}Building application...${NC}"

cd Backend

# Ensure dependencies are downloaded
echo -e "${YELLOW}Downloading dependencies...${NC}"
if ! go mod download 2>&1 | tee /tmp/go_download.log; then
    echo -e "${RED}✗ Failed to download dependencies${NC}"
    cat /tmp/go_download.log
    cd ..
    exit 1
fi

go mod tidy > /dev/null 2>&1

# Build the application
echo -e "${YELLOW}Compiling...${NC}"
if ! go build -o app ./cmd/main.go 2>&1 | tee /tmp/go_build.log; then
    echo -e "${RED}✗ Build failed${NC}"
    echo -e "${YELLOW}Error details:${NC}"
    cat /tmp/go_build.log
    cd ..
    exit 1
fi

echo -e "${GREEN}✓ Application built${NC}"
cd ..

# ============================================
# Install Frontend (if exists)
# ============================================
if [ -d "Frontend" ] && [ -f "Frontend/package.json" ]; then
    echo -e "${BLUE}Installing frontend...${NC}"
    cd Frontend
    
    # Run npm install with visible output
    if npm install 2>&1; then
        echo -e "${GREEN}✓ Frontend ready${NC}"
    else
        echo -e "${RED}✗ Frontend installation failed${NC}"
        cd ..
        exit 1
    fi
    
    cd ..
fi

# ============================================
# Installation Complete
# ============================================
SERVER_IP=$(hostname -I | awk '{print $1}')
[ -z "$SERVER_IP" ] && SERVER_IP=$(ip route get 1 | awk '{print $(NF-2);exit}')

echo -e "\n${GREEN}✓ System successfully installed${NC}"
echo ""
echo -e "${BLUE}Access URL:${NC}"
echo "  http://$SERVER_IP:$APP_PORT"
echo ""