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

# Get server IP (accessible from network)
SERVER_IP=$(hostname -I | awk '{print $1}')
[ -z "$SERVER_IP" ] && SERVER_IP=$(ip route get 1 | awk '{print $(NF-2);exit}')

# ============================================
# Auto-Detect Changes & Build
# ============================================
echo -e "${BLUE}Checking system status...${NC}"

# Function to check if rebuild is needed
check_rebuild() {
    # Check if binary exists
    if [ ! -f "Backend/app" ]; then
        echo -e "${YELLOW}Binary not found, need to build${NC}"
        return 0
    fi
    
    # Get binary modification time
    local binary_mtime=$(stat -c %Y "Backend/app" 2>/dev/null || echo 0)
    
    # Check Go files for changes
    local go_files_changed=0
    for go_file in $(find Backend -name "*.go" -type f 2>/dev/null); do
        local file_mtime=$(stat -c %Y "$go_file" 2>/dev/null)
        if [ $file_mtime -gt $binary_mtime ]; then
            go_files_changed=1
            break
        fi
    done
    
    if [ $go_files_changed -eq 1 ]; then
        return 0
    fi
    
    # Check go.mod/go.sum
    if [ -f "Backend/go.mod" ] && [ "Backend/go.mod" -nt "Backend/app" ]; then
        return 0
    fi
    
    if [ -f "Backend/go.sum" ] && [ "Backend/go.sum" -nt "Backend/app" ]; then
        return 0
    fi
    
    # Check .env file
    if [ "Backend/.env" -nt "Backend/app" ]; then
        return 0
    fi
    
    return 1
}

# Check if rebuild is needed
if check_rebuild; then
    echo -e "${BLUE}Building...${NC}"
    
    # Backup current binary if exists
    if [ -f "Backend/app" ]; then
        cp Backend/app Backend/app.backup 2>/dev/null || true
    fi
    
    # Build the application
    cd Backend
    
    # Ensure dependencies are up to date
    go mod download > /dev/null 2>&1
    go mod tidy > /dev/null 2>&1
    
    # Build with proper path
    if go build -o app ./cmd/main.go 2>&1 | grep -q "error"; then
        echo -e "${RED}✗ Build failed${NC}"
        # Restore backup if exists
        if [ -f "app.backup" ]; then
            mv app.backup app
        fi
        cd ..
        exit 1
    fi
    
    rm -f app.backup 2>/dev/null
    echo -e "${GREEN}✓ Built${NC}"
    cd ..
fi

# ============================================
# Frontend Dependency Check
# ============================================
if [ -d "Frontend" ] && [ -f "Frontend/package.json" ]; then
    # Check if node_modules exists or package.json is newer
    if [ ! -d "Frontend/node_modules" ] || [ "Frontend/package.json" -nt "Frontend/node_modules" ]; then
        echo -e "${BLUE}Installing frontend dependencies...${NC}"
        cd Frontend
        npm install --silent > /dev/null 2>&1
        cd ..
        echo -e "${GREEN}✓ Frontend ready${NC}"
    fi
fi

# ============================================
# Database Connection Test
# ============================================
if ! mysql -u "$DB_USER" -p"$DB_PASS" -D "$DB_NAME" -e "SELECT 1;" 2>/dev/null; then
    echo -e "${RED}✗ Database connection failed${NC}"
    exit 1
fi

# ============================================
# Service Startup
# ============================================
cleanup() {
    echo ""
    
    # Kill backend if running
    if [ ! -z "$BACKEND_PID" ] && kill -0 $BACKEND_PID 2>/dev/null; then
        kill $BACKEND_PID
        wait $BACKEND_PID 2>/dev/null
    fi
    
    # Kill frontend if running
    if [ ! -z "$FRONTEND_PID" ] && kill -0 $FRONTEND_PID 2>/dev/null; then
        kill $FRONTEND_PID
        wait $FRONTEND_PID 2>/dev/null
    fi
    
    exit 0
}

trap cleanup SIGINT SIGTERM EXIT

# Start Backend
cd Backend
./app > backend.log 2>&1 &
BACKEND_PID=$!
cd ..

# Wait for backend to start
sleep 3

# Verify backend is running
if ! ps -p $BACKEND_PID > /dev/null; then
    echo -e "${RED}✗ Backend failed to start${NC}"
    exit 1
fi

# Start Frontend if exists
FRONTEND_PID=""
FRONTEND_PORT=""

if [ -d "Frontend" ] && [ -f "Frontend/package.json" ]; then
    # Find available port (start from APP_PORT + 1)
    for port in $(seq $((APP_PORT + 1)) $((APP_PORT + 10))); do
        if ! lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            FRONTEND_PORT=$port
            break
        fi
    done
    
    if [ -z "$FRONTEND_PORT" ]; then
        FRONTEND_PORT=$((APP_PORT + 1))
    fi
    
    cd Frontend
    
    # Check for development script
    if grep -q '"dev"' package.json; then
        npm run dev -- --host 0.0.0.0 --port $FRONTEND_PORT > frontend.log 2>&1 &
        FRONTEND_PID=$!
    elif grep -q '"start"' package.json; then
        PORT=$FRONTEND_PORT npm start > frontend.log 2>&1 &
        FRONTEND_PID=$!
    fi
    
    cd ..
    sleep 2
fi

# ============================================
# Display System Information
# ============================================
echo ""
echo -e "${GREEN}✓ System started successfully${NC}"
echo ""
echo -e "${BLUE}Backend:${NC} http://$SERVER_IP:$APP_PORT"

if [ ! -z "$FRONTEND_PORT" ] && [ ! -z "$FRONTEND_PID" ]; then
    echo -e "${BLUE}Frontend:${NC} http://$SERVER_IP:$FRONTEND_PORT"
fi

echo ""
echo -e "${YELLOW}Press Ctrl+C to stop${NC}"
echo ""

# ============================================
# Keep running
# ============================================
while true; do
    sleep 10
    
    # Check backend health
    if ! ps -p $BACKEND_PID > /dev/null; then
        echo -e "\n${RED}✗ Backend stopped${NC}"
        break
    fi
done

cleanup