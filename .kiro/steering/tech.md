# Technology Stack

## Backend

- **Language**: Go 1.24+
- **Framework**: Gin (web framework)
- **ORM**: GORM
- **Database**: PostgreSQL 15+
- **Authentication**: JWT
- **Scheduler**: robfig/cron
- **Search**: Meilisearch (optional)
- **Monitoring**: Prometheus metrics

### Key Backend Libraries

- `github.com/gin-gonic/gin` - HTTP web framework
- `gorm.io/gorm` - ORM for database operations
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/robfig/cron/v3` - Scheduled tasks
- `github.com/meilisearch/meilisearch-go` - Search integration
- `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram bot
- `github.com/silenceper/wechat/v2` - WeChat integration
- `google.golang.org/api` - Google API integration

## Frontend

- **Framework**: Nuxt.js 3.8+
- **UI Library**: Naive UI 2.42+
- **Language**: TypeScript 5.0+
- **Styling**: Tailwind CSS 3.x
- **State Management**: Pinia
- **Package Manager**: pnpm 9.13+

### Key Frontend Libraries

- `nuxt` - Vue.js full-stack framework
- `naive-ui` - Vue 3 component library
- `@nuxtjs/tailwindcss` - Tailwind CSS integration
- `@pinia/nuxt` - State management
- `chart.js` - Data visualization
- `qr-code-styling` - QR code generation

## Infrastructure

- **Containerization**: Docker + Docker Compose
- **Reverse Proxy**: Nginx
- **Deployment**: Multi-stage Docker builds

## Common Commands

### Backend Development

```bash
# Run backend server
go run main.go

# Build for production (Linux)
./scripts/build.sh build-linux

# Build for current platform
./scripts/build.sh build

# Check version
./main version

# Run with environment variables
PORT=8080 go run main.go
```

### Frontend Development

```bash
cd web

# Install dependencies
pnpm install

# Development server
pnpm dev

# Build for production
pnpm build

# Preview production build
pnpm preview
```

### Docker

```bash
# Build and start all services
docker-compose up -d

# Build specific service
docker-compose build backend
docker-compose build frontend

# View logs
docker-compose logs -f backend

# Stop all services
docker-compose down
```

### Database

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U postgres -d url_db

# Run migrations (handled automatically on startup)
```

### Build Scripts

```bash
# Build backend binary
./scripts/build.sh build-linux [output-name]

# Build Docker images
./scripts/docker-build.sh build [version]

# Update version
./scripts/version.sh patch|minor|major
```

## Environment Configuration

Key environment variables (see `.env` file):

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Database connection
- `PORT` - Backend server port (default: 8080)
- `TIMEZONE` - Timezone setting (default: Asia/Shanghai)
- `LOG_LEVEL` - Logging level (DEBUG, INFO, WARN, ERROR, FATAL)
- `GIN_MODE` - Gin mode (debug, release, test)
- `NUXT_PUBLIC_API_SERVER` - Frontend API endpoint

## Testing

Currently no automated test suite. Manual testing through:
- API endpoints via curl/Postman
- Frontend UI testing
- Integration testing via Docker Compose
