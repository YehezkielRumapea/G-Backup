#!/bin/bash

echo "=== Backup App Installation ==="

# Check if in correct directory
if [ ! -d "Backend" ]; then
    echo "Error: Run from G-Backup directory"
    exit 1
fi

# Database config
echo ""
echo "Database Configuration:"
read -p "MySQL Username [yehezkiel]: " db_user
db_user=${db_user:-yehezkiel}
read -sp "MySQL Password [kiel]: " db_pass
db_pass=${db_pass:-kiel}
echo ""
read -p "Host [localhost]: " db_host
db_host=${db_host:-localhost}
read -p "Port [3306]: " db_port
db_port=${db_port:-3306}
read -p "Database [tugas_akhir]: " db_name
db_name=${db_name:-tugas_akhir}

# Rclone config
echo ""
echo "Rclone Configuration:"
read -p "Rclone config path [/home/$USER/.config/rclone/rclone.conf]: " rclone_path
rclone_path=${rclone_path:-/home/$USER/.config/rclone/rclone.conf}

# Generate JWT
jwt_secret=$(openssl rand -base64 32 2>/dev/null || echo "secret-$(date +%s)")

# Create .env
cat > Backend/.env << EOF
DB_USER=$db_user
DB_PASS=$db_pass
DB_HOST=$db_host
DB_PORT=$db_port
DB_NAME=$db_name
JWT_SECRET_KEY=$jwt_secret
RCLONE_CONFIG=$rclone_path
EOF

echo "✓ Created Backend/.env"

# Install backend
echo ""
echo "Installing backend..."
cd Backend
go mod download
go build -o backend-app ./cmd/backend_app || exit 1
cd ..
echo "✓ Backend installed"

# Install frontend if exists
if [ -d "frontend" ] && [ -f "frontend/package.json" ]; then
    echo "Installing frontend..."
    cd frontend
    npm install --silent
    cd ..
    echo "✓ Frontend installed"
fi

echo ""
echo "=== Installation Complete ==="
echo "Run: ./start.sh"
echo ""