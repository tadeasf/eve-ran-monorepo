# EVE Ran

![EVE Ran Services Design](https://github.com/tadeasf/eve-ran-monorepo/raw/main/docs/ran-services-design.png)

**EVE Ran** is a web application for tracking and analyzing character kills in **EVE Online**. It provides PVP statistics and performance metrics for members, squads, and sigs inside [Goonswarm Federation](https://goonfleet.com/).

## Table of Contents

- [Features](#features)
- [Demo](#demo)
- [Project Overview](#project-overview)
- [Local Development](#local-development)
- [Deployment](#deployment)
- [Additional Documentation](#additional-documentation)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Character Kill Tracking**: Monitor and analyze killmail data for tracked characters
- **PVP Statistics**: View aggregate statistics, ISK destroyed, and point values
- **Competition Leaderboards**: Year-to-date and monthly rankings by ISK and points
- **Trend Analysis**: Track performance trends over time
- **Region & System Filtering**: Filter kills by EVE Online regions and date ranges
- **Admin Dashboard**: Manage characters, configure competition settings, and view analytics
- **Activity Feed**: Recent kills feed with comments and achievements
- **Real-time Data**: Periodic updates from [zKillboard](https://zkillboard.com/) and [EVE ESI](https://esi.evetech.net/)
- **API Security**: X-API-Key authentication for all endpoints
- **Health Monitoring**: Kubernetes-ready health check endpoints
- **Data Validation**: Comprehensive validation for ESI and zKillboard responses
- **Optimized Performance**: Database indexing and connection pooling

## Demo

![Frontend Demo 1](https://github.com/tadeasf/eve-ran-monorepo/raw/main/docs/ran-frontend-1.png)

![Frontend Demo 2](https://github.com/tadeasf/eve-ran-monorepo/raw/main/docs/ran-frontend-2.png)

## Project Overview

### Architecture

**Stack:**
- Frontend: Next.js 14, React 18, TypeScript, Tailwind CSS, shadcn/ui
- Backend: Go 1.23, Gin framework, GORM
- Database: PostgreSQL
- Reverse Proxy: Caddy

### Backend Core Components

**Services (`backend/src/services/`):**
- `esi.go`: EVE Online ESI API integration with rate limiting
- `esiErrorManager.go`: Global rate limit manager for ESI requests

**Database Models (`backend/src/db/models/`):**
- Character, Kill, Region, System, Constellation, Item models
- Competition settings and tracking

**Jobs (`backend/src/jobs/`):**
- `killCron.go`: Periodic killmail fetching from zKillboard
- `killEnhance.go`: Enrichment with ESI data
- `typesFetcher.go`: Item type metadata updates

**Routes (`backend/src/routes/`):**
- RESTful API endpoints for characters, kills, regions, systems, and competition data

### Frontend Core Components

**Pages (`frontend/src/app/`):**
- `/`: Main dashboard with kill statistics and charts
- `/competition`: Competition leaderboards (YTD and monthly)
- `/admin`: Admin panel for character management and settings

**Key Components (`frontend/src/app/components/`):**
- `CharacterTable.tsx`: Character statistics display
- `FilterControls.tsx`: Region and date filtering
- `TotalKillsChart.tsx` / `TotalIskChart.tsx`: Trend visualization
- `AdminDashboard.tsx`: Admin interface with sidebar navigation
- `AdminCompetitionSettings.tsx`: Competition configuration

**API Client (`frontend/src/app/lib/eveApi.ts`):**
- Centralized API communication with backend
- React Query integration for caching

## Local Development

### Prerequisites

- Go 1.23+
- Node.js 18+
- PostgreSQL 14+
- Docker & Docker Compose (optional)

### Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/tadeasf/eve-ran-monorepo.git
   cd eve-ran-monorepo
   ```

2. **Automated setup (recommended):**
   ```bash
   ./scripts/setup.sh
   ```

   Or manual setup:

3. **Backend setup:**
   ```bash
   cd backend
   cp .env.example .env
   # Generate API key: openssl rand -base64 32
   # Edit .env with your PostgreSQL credentials and API key
   go mod download
   go run src/main.go
   ```

4. **Frontend setup:**
   ```bash
   cd frontend
   cp .env.example .env.local
   # Edit .env.local with your backend API URL and API key (must match backend)
   npm install
   npm run dev
   ```

5. **Database setup:**
   - Create a PostgreSQL database
   - Migrations run automatically on backend startup

### Development URLs

- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080`
- API Documentation: `http://localhost:8080/swagger/index.html`

## Deployment

For production deployment instructions, see **[Deployment Guide](./docs/DEPLOYMENT.md)**.

**Quick Start with Docker Compose:**
```bash
docker-compose up -d
```

This starts all services:
- Frontend on port 12921
- Backend API on port 12922
- PostgreSQL on port 12920

Configure your reverse proxy (Caddy/Nginx) to route traffic to these ports.

## Additional Documentation

- **[Implementation Updates](./docs/IMPLEMENTATION_UPDATES.md)**: Recent improvements and new features (v2.1.0)
- **[ESI Rate Limiting](./docs/ESI_RATE_LIMITING_UPDATE.md)**: Details on EVE Online ESI API rate limiting implementation
- **[Deployment Guide](./docs/DEPLOYMENT.md)**: Production deployment instructions and configuration

## Contributing

Contributions are welcome! Please submit a **Pull Request** with your changes.

For feature requests, open an **Issue** with a description of the desired functionality.

## License

This project is licensed under the GPL-3.0 License.
