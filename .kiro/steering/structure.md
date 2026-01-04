# Project Structure

## Root Directory Layout

```
urlDB/
├── cmd/                    # Command-line tools (e.g., google-index)
├── common/                 # Cloud storage platform integrations
├── config/                 # Configuration management
├── db/                     # Database layer
├── handlers/               # HTTP request handlers (controllers)
├── middleware/             # HTTP middleware (auth, logging, etc.)
├── monitor/                # Monitoring and metrics
├── pkg/                    # Reusable packages (bing, google)
├── scheduler/              # Scheduled tasks (cron jobs)
├── services/               # Business logic services
├── task/                   # Background task processors
├── utils/                  # Utility functions
├── web/                    # Frontend Nuxt.js application
├── scripts/                # Build and deployment scripts
├── migrations/             # Database migration files
├── data/                   # Runtime data (sitemaps, credentials)
├── uploads/                # User uploaded files
├── logs/                   # Application logs
├── main.go                 # Application entry point
├── go.mod, go.sum          # Go dependencies
├── Dockerfile              # Multi-stage Docker build
└── docker-compose.yml      # Docker orchestration
```

## Backend Architecture

### Database Layer (`db/`)

- `entity/` - GORM models (database tables)
- `dto/` - Data Transfer Objects (API request/response)
- `repo/` - Repository pattern implementations
- `converter/` - Entity ↔ DTO converters
- `connection.go` - Database initialization

**Pattern**: Repository pattern with generic base repository

### Handlers (`handlers/`)

HTTP request handlers following REST conventions:
- One handler per resource type
- Dependency injection via `SetRepositoryManager()`
- Standard CRUD operations: Get, Create, Update, Delete
- Response helpers in `response.go`

### Services (`services/`)

Business logic layer:
- `meilisearch_service.go` - Search integration
- `telegram_bot_service.go` - Telegram bot
- `wechat_bot_service.go` - WeChat integration
- Injected into handlers as needed

### Task System (`task/`)

Background task processors:
- `task_processor.go` - Base task processor interface
- `transfer_processor.go` - Resource transfer tasks
- `expansion_processor.go` - Account expansion tasks
- `google_index_processor.go` - Google indexing tasks

### Scheduler (`scheduler/`)

Cron-based scheduled jobs:
- `hot_drama.go` - Fetch hot drama data
- `ready_resource.go` - Process pending resources
- `sitemap.go` - Generate sitemaps
- `google_index.go` - Submit to Google
- `cache_cleaner.go` - Clean expired cache

### Common (`common/`)

Cloud storage platform implementations:
- `base_pan.go` - Base interface for all platforms
- `*_pan.go` - Platform-specific implementations (Quark, Xunlei, Alipan, etc.)
- `pan_factory.go` - Factory pattern for creating platform instances

### Configuration (`config/`)

Unified configuration management:
- `config.go` - ConfigManager with in-memory cache
- `global.go` - Global config instance
- `sync.go` - Config synchronization
- Supports database-backed config with environment variable fallback

## Frontend Architecture (`web/`)

### Nuxt.js Structure

```
web/
├── assets/                 # Static assets (CSS, images)
├── components/             # Vue components
│   ├── Admin/             # Admin-specific components
│   ├── User/              # User-specific components
│   └── QRCode/            # QR code components
├── composables/            # Vue composables (hooks)
├── config/                 # Frontend configuration
├── layouts/                # Page layouts
├── middleware/             # Route middleware
├── pages/                  # File-based routing
│   ├── admin/             # Admin pages
│   ├── user/              # User pages
│   └── r/                 # Resource detail pages
├── plugins/                # Nuxt plugins
├── public/                 # Public static files
├── server/                 # Server-side API routes
├── stores/                 # Pinia state stores
├── nuxt.config.ts          # Nuxt configuration
├── tailwind.config.js      # Tailwind configuration
└── package.json            # Dependencies
```

### Component Organization

- **Admin components**: Prefixed with `Admin/` for admin panel
- **User components**: Prefixed with `User/` for user dashboard
- **Shared components**: Root level (SearchModal, QrCodeModal, etc.)

### State Management

Pinia stores in `stores/`:
- `user.ts` - User authentication state
- `resource.ts` - Resource data
- `systemConfig.ts` - System configuration
- `task.ts` - Task management

### API Integration

- `composables/useApi.ts` - Base API client
- `composables/useApiFetch.ts` - Nuxt-specific fetch wrapper
- Server-side API routes in `server/api/` for SSR

## Key Patterns

### Repository Pattern

All database access goes through repositories:
```go
type BaseRepository[T any] interface {
    Create(entity *T) error
    FindByID(id uint) (*T, error)
    FindAll() ([]T, error
    Update(entity *T) error
    Delete(id uint) error
}
```

### Dependency Injection

Services and repositories injected via setter functions:
```go
handlers.SetRepositoryManager(repoManager)
services.SetMeilisearchManager(meilisearchManager)
```

### Task Processor Pattern

Background tasks implement `TaskProcessor` interface:
```go
type TaskProcessor interface {
    GetType() string
    Process(task *entity.Task) error
}
```

### Configuration Management

Centralized config with caching:
```go
configManager.GetConfigBool(entity.ConfigKeyAutoTransferEnabled)
configManager.SetConfig(key, value)
```

## File Naming Conventions

- **Go files**: snake_case (e.g., `resource_handler.go`)
- **Vue files**: PascalCase (e.g., `SearchModal.vue`)
- **TypeScript files**: camelCase (e.g., `useApi.ts`)
- **Test files**: `*_test.go` (Go convention)

## Import Organization

Go imports grouped in order:
1. Standard library
2. External packages
3. Internal packages (github.com/ctwj/urldb/...)
