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
echo "│        Starting Backup System               │"
echo "└──────────────────────────────────────────────┘"
echo -e "${NC}"

# ============================================
# Configuration Check
# ============================================
if [ ! -f "Backend/.env" ]; then
    echo -e "${RED}Configuration not found${NC}"
    echo "Please run ./install.sh first"
    exit 1
fi

# Load configuration
source <(grep -v '^#' Backend/.env | sed 's/^/export /')

# Get server IP
SERVER_IP=$(hostname -I | awk '{print $1}')
[ -z "$SERVER_IP" ] && SERVER_IP="127.0.0.1"

# ============================================
# Auto-Detect Changes & Build
# ============================================
echo -e "${BLUE}Checking system status...${NC}"

# Function to check if rebuild is needed
check_rebuild() {
    local rebuild_needed=1
    
    # Check if binary exists
    if [ ! -f "Backend/app" ]; then
        echo -e "${YELLOW}Binary not found, need to build${NC}"
        return 0
    fi
    
    # Get binary modification time
    local binary_mtime=$(stat -c %Y "Backend/app" 2>/dev/null || echo 0)
    
    # Check Go files for changes
    local go_files_changed=0
    for go_file in $(find Backend -name "*.go" -type f); do
        local file_mtime=$(stat -c %Y "$go_file" 2>/dev/null)
        if [ $file_mtime -gt $binary_mtime ]; then
            echo -e "${YELLOW}Go file changed: $(basename "$go_file")${NC}"
            go_files_changed=1
            break
        fi
    done
    
    if [ $go_files_changed -eq 1 ]; then
        return 0
    fi
    
    # Check go.mod
    if [ -f "Backend/go.mod" ] && [ "Backend/go.mod" -nt "Backend/app" ]; then
        echo -e "${YELLOW}Dependencies changed${NC}"
        return 0
    fi
    
    # Check .env file
    if [ "Backend/.env" -nt "Backend/app" ]; then
        echo -e "${YELLOW}Configuration changed${NC}"
        return 0
    fi
    
    echo -e "${GREEN}✓ Code is up to date${NC}"
    return 1
}

# Check if rebuild is needed
if check_rebuild; then
    echo -e "\n${BLUE}Building application...${NC}"
    
    # Backup current binary if exists
    if [ -f "Backend/app" ]; then
        echo -e "${YELLOW}Backing up current binary...${NC}"
        cp Backend/app Backend/app.backup 2>/dev/null || true
    fi
    
    # Build the application
    cd Backend
    if go build -o app ./cmd/backend_app; then
        echo -e "${GREEN}✓ Build successful${NC}"
        # Remove backup if build succeeded
        rm -f app.backup 2>/dev/null
    else
        echo -e "${RED}✗ Build failed${NC}"
        # Restore backup if exists
        if [ -f "app.backup" ]; then
            echo -e "${YELLOW}Restoring previous version...${NC}"
            mv app.backup app
        fi
        exit 1
    fi
    cd ..
else
    echo -e "${GREEN}✓ No rebuild needed${NC}"
fi

# ============================================
# Frontend Dependency Check
# ============================================
if [ -d "Frontend" ] && [ -f "Frontend/package.json" ]; then
    echo -e "\n${BLUE}Checking frontend...${NC}"
    
    # Check if node_modules exists or package.json is newer
    if [ ! -d "Frontend/node_modules" ] || [ "Frontend/package.json" -nt "Frontend/node_modules" ]; then
        echo -e "${YELLOW}Installing/updating frontend dependencies...${NC}"
        cd Frontend
        npm install --silent
        cd ..
        echo -e "${GREEN}✓ Frontend dependencies updated${NC}"
    else
        echo -e "${GREEN}✓ Frontend dependencies are current${NC}"
    fi
fi

# ============================================
# Service Startup
# ============================================
cleanup() {
    echo -e "\n${YELLOW}Shutting down services...${NC}"
    
    # Kill backend if running
    if [ ! -z "$BACKEND_PID" ] && kill -0 $BACKEND_PID 2>/dev/null; then
        kill $BACKEND_PID
        echo -e "${GREEN}✓ Backend stopped${NC}"
    fi
    
    # Kill frontend if running
    if [ ! -z "$FRONTEND_PID" ] && kill -0 $FRONTEND_PID 2>/dev/null; then
        kill $FRONTEND_PID
        echo -e "${GREEN}✓ Frontend stopped${NC}"
    fi
    
    exit 0
}

trap cleanup SIGINT SIGTERM

# Start Backend
echo -e "\n${BLUE}Starting backend service...${NC}"
cd Backend
./app &
BACKEND_PID=$!
cd ..

# Wait a moment for backend to start
sleep 2

# Verify backend is running
if ps -p $BACKEND_PID > /dev/null; then
    echo -e "${GREEN}✓ Backend running (PID: $BACKEND_PID)${NC}"
else
    echo -e "${RED}✗ Backend failed to start${NC}"
    exit 1
fi

# Start Frontend if exists
if [ -d "Frontend" ] && [ -f "Frontend/package.json" ]; then
    echo -e "\n${BLUE}Starting frontend service...${NC}"
    
    # Find available port
    FRONTEND_PORT=""
    for port in $(seq $APP_PORT $((APP_PORT + 5))); do
        if ! lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            FRONTEND_PORT=$port
            break
        fi
    done
    
    FRONTEND_PORT=${FRONTEND_PORT:-$((APP_PORT + 1))}
    
    cd Frontend
    
    # Check for development script
    if grep -q '"dev"' package.json; then
        npm run dev -- --host 0.0.0.0 --port $FRONTEND_PORT > frontend.log 2>&1 &
    else
        # Try to serve static files if dist exists
        if [ -d "dist" ]; then
            npx serve -s dist -l $FRONTEND_PORT > frontend.log 2>&1 &
        else
            echo -e "${YELLOW}⚠ No start method found for frontend${NC}"
            FRONTEND_PID=""
        fi
    fi
    
    FRONTEND_PID=$!
    cd ..
    
    sleep 3
    
    if [ ! -z "$FRONTEND_PID" ] && ps -p $FRONTEND_PID > /dev/null; then
        echo -e "${GREEN}✓ Frontend running on port $FRONTEND_PORT${NC}"
    fi
fi

# ============================================
# Display System Information
# ============================================
echo -e "\n${BLUE}══════════════════════════════════════════════${NC}"
echo -e "${GREEN}          System Started Successfully        ${NC}"
echo -e "${BLUE}══════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}Server:${NC} $SERVER_IP"
echo ""
echo -e "${GREEN}Backend API${NC}"
echo -e "  http://$SERVER_IP:$APP_PORT"
echo ""

if [ ! -z "$FRONTEND_PORT" ] && [ ! -z "$FRONTEND_PID" ]; then
    echo -e "${GREEN}Frontend Interface${NC}"
    echo -e "  http://$SERVER_IP:$FRONTEND_PORT"
    echo ""
fi

echo -e "${YELLOW}Process Information${NC}"
echo -e "  Backend PID: $BACKEND_PID"
[ ! -z "$FRONTEND_PID" ] && echo -e "  Frontend PID: $FRONTEND_PID"
echo ""
echo -e "${BLUE}══════════════════════════════════════════════${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo ""

# ============================================
# Health Monitoring
# ============================================
echo -e "${BLUE}Monitoring services...${NC}"

while true; do
    # Check backend health
    if ! ps -p $BACKEND_PID > /dev/null; then
        echo -e "\n${RED}Backend process has stopped${NC}"
        break
    fi
    
    # Check frontend health if running
    if [ ! -z "$FRONTEND_PID" ] && ! ps -p $FRONTEND_PID > /dev/null; then
        echo -e "\n${YELLOW}Frontend process has stopped${NC}"
        FRONTEND_PID=""
    fi
    
    # Check service health via API (optional)
    if curl -s http://localhost:$APP_PORT/health > /dev/null 2>&1; then
        echo -n "."
    else
        echo -e "\n${YELLOW}API health check failed${NC}"
    fi
    
    sleep 10
done

cleanup