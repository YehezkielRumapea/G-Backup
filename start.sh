#!/bin/bash

echo "=== Starting Backup Application ==="

# Cleanup function
cleanup() {
    echo -e "\nStopping application..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    exit 0
}
trap cleanup SIGINT

# Check .env
if [ ! -f "Backend/.env" ]; then
    echo "Error: Backend/.env not found"
    echo "Run ./install.sh first"
    exit 1
fi

# Build backend if needed
if [ ! -f "Backend/backend-app" ]; then
    echo "Building backend..."
    cd Backend
    go build -o backend-app ./cmd/backend_app || exit 1
    cd ..
fi

# Load environment
export $(grep -v '^#' Backend/.env | xargs)

# Start backend
echo "Starting backend..."
cd Backend
./backend-app &
BACKEND_PID=$!
cd ..

sleep 2
echo "✓ Backend started"

# Start frontend if exists
if [ -d "Frontend" ] && [ -f "Frontend/package.json" ]; then
    echo "Starting frontend..."
    cd Frontend
    
    # Check if port 3000 is available
    if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null; then
        PORT=3001
        echo "Port 3000 busy, using 3001"
        npm run dev -- --port 3001 > /dev/null 2>&1 &
    else
        PORT=3000
        npm run dev > /dev/null 2>&1 &
    fi
    
    FRONTEND_PID=$!
    cd ..
    
    sleep 3
    echo "✓ Frontend started"
    
    # Show frontend URL at the end
    echo ""
    echo "========================================"
    echo "Frontend URL: http://localhost:$PORT"
    echo "Backend API:  http://localhost:8080"
    echo "========================================"
else
    echo ""
    echo "========================================"
    echo "Frontend not found"
    echo "Backend API: http://localhost:8080"
    echo "========================================"
fi

echo ""
echo "Press Ctrl+C to stop"
echo ""

wait $BACKEND_PID