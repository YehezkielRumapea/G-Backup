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

# Progress bar function
show_progress() {
    local duration=$1
    local message=$2
    local steps=20
    local step_duration=$(echo "scale=2; $duration / $steps" | bc)
    
    echo -n "$message "
    for i in $(seq 1 $steps); do
        echo -n "▓"
        sleep $step_duration
    done
    echo " Done"
}

install_dep() {
    if ! command -v $1 &>/dev/null; then
        echo -e "${YELLOW}Installing $1...${NC}"
        (sudo apt update -qq && sudo apt install -y $2 -qq) &
        show_progress 2 "Downloading"
        wait
        echo -e "${GREEN}✓ $1 installed${NC}"
    else
        echo -e "${GREEN}✓ $1 available${NC}"
    fi
}

install_dep "go" "golang-go"
install_dep "node" "nodejs"

# NPM needs special handling - check if it actually works
echo -n "Checking npm... "
if command -v npm &>/dev/null && npm --version &>/dev/null; then
    NPM_VERSION=$(npm --version 2>/dev/null)
    echo -e "${GREEN}✓ npm v${NPM_VERSION}${NC}"
else
    echo -e "${YELLOW}not found/broken${NC}"
    echo -e "${YELLOW}Installing npm...${NC}"
    (sudo apt update -qq && sudo apt install -y npm nodejs -qq) &
    show_progress 3 "Downloading"
    wait
    
    # Verify installation
    if npm --version &>/dev/null; then
        NPM_VERSION=$(npm --version 2>/dev/null)
        echo -e "${GREEN}✓ npm v${NPM_VERSION} installed${NC}"
    else
        echo -e "${RED}✗ npm installation failed${NC}"
        echo -e "${YELLOW}Try manually: sudo apt install nodejs npm${NC}"
        exit 1
    fi
fi

install_dep "mysql" "default-mysql-client"

# Rclone
if ! command -v rclone &>/dev/null; then
    echo -e "${YELLOW}Installing rclone...${NC}"
    (curl -fsSL https://rclone.org/install.sh | sudo bash > /dev/null 2>&1) &
    show_progress 4 "Downloading"
    wait
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

# Auto-detect rclone config path based on current user
CURRENT_USER=$(whoami)
RCLONE_PATH="/home/$CURRENT_USER/.config/rclone/rclone.conf"

# Create rclone config directory if not exists
RCLONE_DIR=$(dirname "$RCLONE_PATH")
if [ ! -d "$RCLONE_DIR" ]; then
    mkdir -p "$RCLONE_DIR"
    echo -e "${YELLOW}Created rclone config directory${NC}"
fi

echo -e "${GREEN}✓ Rclone path: $RCLONE_PATH${NC}"

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
echo -n "Downloading dependencies... "
(go mod download 2>&1 > /tmp/go_download.log) &
pid=$!
while kill -0 $pid 2>/dev/null; do
    echo -n "▓"
    sleep 0.3
done
wait $pid
if [ $? -ne 0 ]; then
    echo ""
    echo -e "${RED}✗ Failed to download dependencies${NC}"
    cat /tmp/go_download.log
    cd ..
    exit 1
fi
echo " Done"

go mod tidy > /dev/null 2>&1

# Build the application
echo -n "Compiling... "
(go build -o app ./cmd/main.go 2>&1 > /tmp/go_build.log) &
pid=$!
while kill -0 $pid 2>/dev/null; do
    echo -n "▓"
    sleep 0.2
done
wait $pid
if [ $? -ne 0 ]; then
    echo ""
    echo -e "${RED}✗ Build failed${NC}"
    echo -e "${YELLOW}Error details:${NC}"
    cat /tmp/go_build.log
    cd ..
    exit 1
fi
echo " Done"

echo -e "${GREEN}✓ Application built${NC}"
cd ..

# ============================================
# Install Frontend (if exists)
# ============================================
if [ -d "Frontend" ] && [ -f "Frontend/package.json" ]; then
    echo -e "${BLUE}Installing frontend...${NC}"
    
    # Double-check npm is actually working
    if ! npm --version &>/dev/null; then
        echo -e "${RED}✗ npm is not working${NC}"
        echo -e "${YELLOW}Installing npm...${NC}"
        (sudo apt update -qq && sudo apt install -y npm nodejs -qq) &
        show_progress 3 "Downloading"
        wait
        
        # Final check
        if ! npm --version &>/dev/null; then
            echo -e "${RED}✗ Failed to install npm${NC}"
            echo -e "${YELLOW}Run manually: sudo apt install nodejs npm${NC}"
            exit 1
        fi
    fi
    
    cd Frontend
    
    # Run npm install with progress
    echo -n "Installing packages... "
    (npm install --silent 2>&1 > /tmp/npm_install.log) &
    pid=$!
    while kill -0 $pid 2>/dev/null; do
        echo -n "▓"
        sleep 0.5
    done
    wait $pid
    if [ $? -ne 0 ]; then
        echo ""
        echo -e "${RED}✗ Frontend installation failed${NC}"
        cat /tmp/npm_install.log
        cd ..
        exit 1
    fi
    echo " Done"
    
    echo -e "${GREEN}✓ Frontend ready${NC}"
    cd ..
fi

# ============================================
# Installation Complete
# ============================================

# Get the correct network IP (not NAT/loopback)
get_network_ip() {
    # Get all IPs excluding loopback
    local all_ips=$(ip -4 addr show | grep "inet " | grep -v "127.0.0.1" | \
                    awk '{print $2}' | cut -d/ -f1)
    
    # Get IPs excluding common NAT ranges
    local good_ips=$(echo "$all_ips" | grep -v "^10\.0\.2\." | grep -v "^172\.17\." | grep -v "^169\.254\.")
    
    # Check if we have good IPs
    if [ ! -z "$good_ips" ]; then
        local ip_count=$(echo "$good_ips" | wc -l)
        
        # If only one good IP, use it
        if [ $ip_count -eq 1 ]; then
            echo "$good_ips"
            return 0
        fi
        
        # Multiple good IPs - use priority logic
        # Priority 1: IP from default route interface
        local default_iface=$(ip route | grep "^default" | awk '{print $5}' | head -n1)
        if [ ! -z "$default_iface" ]; then
            local default_ip=$(ip addr show "$default_iface" | grep "inet " | grep -v "127.0.0.1" | \
                              awk '{print $2}' | cut -d/ -f1 | head -n1)
            if [ ! -z "$default_ip" ] && [[ ! $default_ip =~ ^10\.0\.2\. ]]; then
                echo "$default_ip"
                return 0
            fi
        fi
        
        # Priority 2: Prefer 192.168.x.x or 10.x.x.x
        for ip in $good_ips; do
            if [[ $ip =~ ^192\.168\. ]] || [[ $ip =~ ^10\.[0-9]+\.[0-9]+\. ]]; then
                echo "$ip"
                return 0
            fi
        done
        
        # Use first good IP
        echo "$good_ips" | head -n1
        return 0
    fi
    
    # No good IPs found, use NAT IP as fallback
    if [ ! -z "$all_ips" ]; then
        echo "$all_ips" | head -n1
        return 1  # Return 1 to indicate NAT IP is being used
    fi
    
    # Absolute fallback
    echo "127.0.0.1"
    return 2
}

SERVER_IP=$(get_network_ip)
IP_STATUS=$?

# Show all available IPs
ALL_IPS=$(ip -4 addr show | grep "inet " | grep -v "127.0.0.1" | \
          awk '{print $2}' | cut -d/ -f1)
IP_COUNT=$(echo "$ALL_IPS" | wc -l)

echo -e "\n${GREEN}✓ System successfully installed${NC}"
echo ""

# Check if using NAT IP (status = 1)
if [ $IP_STATUS -eq 1 ]; then
    echo -e "${BLUE}Access URL:${NC}"
    echo "  http://$SERVER_IP:$APP_PORT"
    echo ""
    echo -e "${YELLOW}⚠ Network Notice:${NC}"
    echo "  Using NAT IP - may not be accessible from other devices"
    echo "  Fix: VirtualBox → Settings → Network → Bridged Adapter"
    
elif [ $IP_COUNT -gt 1 ]; then
    echo -e "${BLUE}Access URL (multiple networks):${NC}"
    while IFS= read -r ip; do
        if [ "$ip" = "$SERVER_IP" ]; then
            echo -e "  ${GREEN}➜ http://$ip:$APP_PORT${NC} (primary)"
        else
            echo -e "    http://$ip:$APP_PORT"
        fi
    done <<< "$ALL_IPS"
else
    echo -e "${BLUE}Access URL:${NC}"
    echo "  http://$SERVER_IP:$APP_PORT"
fi
echo ""