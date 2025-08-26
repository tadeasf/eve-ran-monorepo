# EVE Ran Deployment Guide

## Overview
This application consists of three main components:
- Frontend (Next.js) - runs on port 12921
- Backend API (Go) - runs on port 12922  
- PostgreSQL Database - runs on port 12920

## Environment Setup

### Backend Environment (.env)
Create `backend/.env` with:
```
DB_HOST=postgres
DB_PORT=5432
DB_USER=eve
DB_PASSWORD=your_secure_password
DB_NAME=eve
```

### Frontend Environment (.env.local)
The frontend is configured to use `https://api.tundragon.space` as the API URL.

## Deployment

### Using Docker Compose
```bash
docker-compose up -d
```

This will start:
- Frontend on http://localhost:12921
- Backend API on http://localhost:12922
- PostgreSQL on localhost:12920

### Machine-wide Reverse Proxy Setup
Configure your machine-wide Caddyfile to include:

```
tundragon.space {
    reverse_proxy localhost:12921
}

api.tundragon.space {
    reverse_proxy localhost:12922
}
```

### SSL/TLS
The machine-wide Caddy instance will handle SSL certificate management automatically.

## Architecture Changes
- Removed Caddy from Docker Compose
- Direct port exposure for manual reverse proxy setup
- Simplified container dependencies
- Ready for production deployment

## Ports
- 12920: PostgreSQL (internal access)
- 12921: Frontend (external access via reverse proxy)
- 12922: Backend API (external access via reverse proxy)
