#!/bin/bash

# EVE Ran Monorepo - Installation & Update Script
# This script installs dependencies and builds both frontend and backend

set -e  # Exit on error

echo "========================================="
echo "EVE Ran Monorepo - Setup Script"
echo "========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if we're in the correct directory
if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}Error: docker-compose.yml not found. Please run this script from the monorepo root.${NC}"
    exit 1
fi

# Backend setup
echo -e "${YELLOW}Setting up backend...${NC}"
cd backend

# Check if .env exists
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}Creating .env from .env.example...${NC}"
    cp .env.example .env
    echo -e "${RED}⚠️  Please update backend/.env with your configuration!${NC}"
fi

# Build backend
echo -e "${GREEN}Building Go backend...${NC}"
go mod download
go build -o main src/main.go

echo -e "${GREEN}✓ Backend setup complete!${NC}"
cd ..

# Frontend setup
echo -e "${YELLOW}Setting up frontend...${NC}"
cd frontend

# Check if .env.local exists
if [ ! -f ".env.local" ]; then
    echo -e "${YELLOW}Creating .env.local from .env.example...${NC}"
    cp .env.example .env.local
    echo -e "${RED}⚠️  Please update frontend/.env.local with your configuration!${NC}"
fi

# Install dependencies
echo -e "${GREEN}Installing frontend dependencies...${NC}"
npm install

# Build frontend
echo -e "${GREEN}Building Next.js frontend...${NC}"
npm run build

echo -e "${GREEN}✓ Frontend setup complete!${NC}"
cd ..

echo ""
echo -e "${GREEN}========================================="
echo "✓ Setup Complete!"
echo "=========================================${NC}"
echo ""
echo "Next steps:"
echo "1. Update backend/.env with your database credentials and API key"
echo "2. Update frontend/.env.local with your API URL and key"
echo "3. Generate a secure API key:"
echo "   openssl rand -base64 32"
echo "4. Start the services:"
echo "   docker-compose up -d"
echo ""
echo "Endpoints:"
echo "  - Frontend: http://localhost:12921"
echo "  - Backend API: http://localhost:12922"
echo "  - API Documentation: http://localhost:12922/swagger/index.html"
echo "  - Health Check: http://localhost:12922/health"
echo ""
