# EVE Ran Deployment Guide

## Components

- **Frontend** (Next.js): Port 12921
- **Backend** (Go): Port 12922
- **Database** (PostgreSQL): Port 12920

## Environment Configuration

### Backend (`backend/.env`)
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=eve
DB_PASSWORD=your_secure_password
DB_NAME=eve
```

### Frontend (`frontend/.env.local`)
Configure the backend API URL:
```env
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
```

## Deployment Steps

### 1. Docker Compose Deployment

```bash
docker-compose up -d
```

This starts all services with exposed ports:
- Frontend: `http://localhost:12921`
- Backend: `http://localhost:12922`
- PostgreSQL: `localhost:12920` (internal only)

### 2. Reverse Proxy Configuration

**Caddy Example:**
```caddy
yourdomain.com {
    reverse_proxy localhost:12921
}

api.yourdomain.com {
    reverse_proxy localhost:12922
}
```

**Nginx Example:**
```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;
    
    location / {
        proxy_pass http://localhost:12921;
    }
}

server {
    listen 443 ssl;
    server_name api.yourdomain.com;
    
    location / {
        proxy_pass http://localhost:12922;
    }
}
```

### 3. SSL/TLS

- **Caddy**: Automatic certificate management
- **Nginx**: Use Certbot with Let's Encrypt

## Architecture Notes

- Direct port exposure (no Caddy in Docker)
- Machine-wide reverse proxy handles SSL
- Simplified container dependencies
- Production-ready configuration
