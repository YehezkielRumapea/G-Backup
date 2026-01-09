#!/bin/bash

set -e

# ==============================
# Terminal Colors
# ==============================
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}"
echo "┌──────────────────────────────────────────────┐"
echo "│        Backup Management System              │"
echo "│              Service Startup                 │"
echo "└──────────────────────────────────────────────┘"
echo -e "${NC}"

# ==============================
# Utility Functions
# ==============================
get_ip() {
    hostname -I | awk '{print $1}'
}

is_port_available() {
    local port=$1
    ! lsof -Pi :"$port" -sTCP:LISTEN -t >/dev/null 2>&1
}

select_application_port() {
    local ports=(5173 3000 3001 8080)

    # 1. Prioritaskan port dari .env jika tersedia & bebas
    if [ -n "$APP_PORT" ] && is_port_available "$APP_PORT"; then
        echo "$APP_PORT"
        return
    fi

    # 2. Fallback ke daftar port standar
    for port in "${ports[@]}"; do
        if is_port_available "$port"; then
            echo "$port"
            return
        fi
    done

    # 3. Fallback terakhir (hard default)
    echo "5173"
}

# ==============================
# Environment Validation
# ==============================
if [ ! -f "Backend/.env" ]; then
    echo -e "${RED}✗ Configuration file not found (Backend/.env)${NC}"
    echo -e "${YELLOW}Please run install.sh before starting the system.${NC}"
    exit 1
fi

echo -e "${BLUE}Loading system configuration...${NC}"
export $(grep -v '^#' Backend/.env | xargs)

# ==============================
# Port Resolution
# ==============================
echo -e "${BLUE}Resolving application port...${NC}"
FINAL_APP_PORT=$(select_application_port)

if [ "$FINAL_APP_PORT" != "$APP_PORT" ]; then
    echo -e "${YELLOW}Requested port unavailable. Using fallback port: $FINAL_APP_PORT${NC}"
fi

export APP_PORT=$FINAL_APP_PORT

# ==============================
# Service Startup
# ==============================
echo ""
echo -e "${BLUE}Starting application service...${NC}"

cd Backend

nohup ./app > system.log 2>&1 &
SERVICE_PID=$!

cd ..

sleep 1

if ps -p $SERVICE_PID > /dev/null; then
    echo -e "${GREEN}✓ Service started successfully${NC}"
else
    echo -e "${RED}✗ Failed to start service${NC}"
    exit 1
fi

# ==============================
# Access Information
# ==============================
SERVER_IP=$(get_ip)

echo ""
echo -e "${BLUE}┌──────────────────────────────────────────────┐"
echo -e "│               Service Status                 │"
echo -e "└──────────────────────────────────────────────┘${NC}"
echo ""
echo -e "${GREEN}System is now running.${NC}"
echo ""
echo -e "${YELLOW}Access URL:${NC}"
echo -e "  http://${SERVER_IP}:${APP_PORT}"
echo ""
echo -e "${BLUE}Process ID:${NC} ${SERVICE_PID}"
echo ""
echo -e "${BLUE}Logs:${NC}"
echo "  Backend/system.log"
echo ""
