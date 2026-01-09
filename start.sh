#!/bin/bash

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}"
echo "┌──────────────────────────────────────────────┐"
echo "│         Starting Backup Server              │"
echo "└──────────────────────────────────────────────┘"
echo -e "${NC}"

# Check if .env exists
if [ ! -f "Backend/.env" ]; then
    echo -e "${RED}Error: Configuration not found${NC}"
    echo "Run ./install.sh first"
    exit 1
fi

# Get server IP
get_ip() {
    ip=$(hostname -I | awk '{print $1}')
    [ -z "$ip" ] && ip=$(ip addr show | grep -oP 'inet \K[\d.]+' | grep -v 127.0.0.1 | head -n1)
    echo $ip
}

SERVER_IP=$(get_ip)

# Load configuration
source <(grep -v '^#' Backend/.env | sed 's/^/export /')

# Build backend if needed
if [ ! -f "Backend/backend-app" ]; then
    echo -e "${YELLOW}Building backend...${NC}"
    cd Backend
    go build -o backend-app ./cmd/backend_app || exit 1
    cd ..
fi

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}Stopping server...${NC}"
    [ ! -z "$BACKEND_PID" ] && kill $BACKEND_PID 2>/dev/null
    [ ! -z "$FRONTEND_PID" ] && kill $FRONTEND_PID 2>/dev/null
    exit 0
}

trap cleanup SIGINT

# Start backend
echo -e "${BLUE}Starting backend...${NC}"
cd Backend
SERVER_HOST=0.0.0.0 ./backend-app &
BACKEND_PID=$!
cd ..

sleep 2
echo -e "${GREEN}✓ Backend running on port ${BACKEND_PORT}${NC}"

# Start frontend if exists
if [ -d "frontend" ] && [ -f "frontend/package.json" ]; then
    echo -e "${BLUE}Starting frontend...${NC}"
    cd frontend
    
    # Find available port
    for port in {$FRONTEND_PORT,3001,3002,3003}; do
        if ! lsof -Pi :$port -sTCP:LISTEN -t >/dev/null; then
            ACTUAL_PORT=$port
            break
        fi
    done
    
    # Serve with host binding
    npm run dev -- --host 0.0.0.0 --port $ACTUAL_PORT > /dev/null 2>&1 &
    FRONTEND_PID=$!
    cd ..
    
    sleep 3
    echo -e "${GREEN}✓ Frontend running on port ${ACTUAL_PORT}${NC}"
fi

echo ""
echo -e "${BLUE}══════════════════════════════════════════════${NC}"
echo -e "${GREEN}           Server Access Information          ${NC}"
echo -e "${BLUE}══════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}Server IP:${NC} ${SERVER_IP}"
echo ""

if [ ! -z "$ACTUAL_PORT" ]; then
    echo -e "${GREEN}Frontend URL:${NC}"
    echo -e "  http://${SERVER_IP}:${ACTUAL_PORT}"
    echo ""
fi

echo -e "${GREEN}Backend API:${NC}"
echo -e "  http://${SERVER_IP}:${BACKEND_PORT}"
echo ""
echo -e "${BLUE}══════════════════════════════════════════════${NC}"
echo -e "${YELLOW}From other devices, use the URLs above.${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop the server${NC}"
echo ""

# Keep running
wait $BACKEND_PID 2>/dev/null