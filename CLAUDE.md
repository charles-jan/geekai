# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GeekAI is a full-stack AI assistant platform featuring a Go backend API and Vue 3 frontend. It integrates multiple AI services including ChatGPT, Claude, MidJourney, Stable Diffusion, DALL-E, Suno (music), and video generation. The project includes both user-facing applications and an administrative management system.

## Repository Structure

```
/
├── api/          # Go backend (Gin + MySQL + Redis + fx)
├── web/          # Vue 3 frontend (Element Plus + Vant)
├── database/     # MySQL schema files and migrations
├── deploy/       # Docker Compose deployment configuration
├── build/        # Docker build scripts
└── desktop/      # Desktop client (Electron-based)
```

## Development Commands

### Backend (api/)

```bash
# Run the backend server
cd api && go run main.go

# Build for different architectures
cd api && make amd64    # Linux AMD64
cd api && make arm64    # Linux ARM64
cd api && make          # Build both

# Run tests
cd api && go test ./test/...
cd api && go test ./test/service/captcha_service_test.go -v

# Configuration
cp api/config.sample.toml api/config.toml
# Edit api/config.toml with your settings
```

### Frontend (web/)

```bash
# Install dependencies
cd web && npm install

# Start development server (port 8888)
cd web && npm run dev

# Build for production
cd web && npm run build

# Lint code
cd web && npm run lint
```

### Docker Deployment

```bash
# Start all services (MySQL, Redis, API, Web)
cd deploy && docker-compose up -d

# Build Docker images
cd build && ./build.sh <version> <amd64|arm64> [push]
# Example: ./build.sh v4.1.10 amd64 push
```

## Architecture Overview

### Backend Architecture (api/)

**Technology Stack:**
- **Framework**: Gin HTTP framework
- **Dependency Injection**: uber-go/fx
- **ORM**: GORM
- **Databases**: MySQL (primary), Redis (cache/sessions), LevelDB (embedded storage)
- **Real-time**: WebSocket support for chat streaming

**Layered Architecture:**
- `main.go` - Application entry point with fx dependency injection wiring
- `core/` - Core server, middleware, configuration, type definitions
- `handler/` - HTTP request handlers (controllers), split into public and `admin/` endpoints
- `service/` - Business logic services (dalle, mj, sd, suno, video, payment, oss, sms)
- `store/` - Data persistence layer with `model/` (database models) and `vo/` (value objects)
- `utils/` - Utility functions and response helpers

**Key Patterns:**
- All dependencies registered via `fx.Provide()` in main.go
- JWT token authentication with separate flows for users and admins
- Redis-based session management
- WebSocket connections with token-based authentication
- Background job scheduling via XXL-Job

### Frontend Architecture (web/)

**Technology Stack:**
- **Framework**: Vue 3 Composition API
- **UI Libraries**: Element Plus (desktop), Vant (mobile)
- **State Management**: Pinia stores
- **Styling**: Stylus + Tailwind CSS
- **Build Tool**: Vue CLI 5

**Dual Interface Design:**
- Desktop web interface at `/home/*`, `/chat`, `/mj`, etc.
- Mobile-optimized interface at `/mobile/*`
- Shared components in `src/components/`, mobile-specific in `src/components/mobile/`

**Key Pinia Stores:**
- `sharedata.js` - Global app state, WebSocket connection, theme management
- `session.js` - User authentication and token management
- `cache.js` - Client-side session and system info caching
- `theme.js` - Theme switching functionality
- `sidebar.js` - Sidebar state management

**API Communication:**
- Axios-based HTTP client in `src/utils/http.js`
- Automatic token injection for authentication
- WebSocket auto-reconnection in App.vue
- Environment-based API host configuration

## Database Schema

- **Database Files**: `database/` contains versioned SQL schema files
- **Initial Schema**: Use `database/geekai_plus-v4.1.8.sql` (latest)
- **Migrations**: Update scripts in `database/update-v*.sql` format
- **Connection**: Configured in `api/config.toml` under MysqlDns

## Configuration

### Backend Configuration (api/config.toml)

Key configuration sections:
- `Listen` - Server bind address (default: 0.0.0.0:5678)
- `MysqlDns` - MySQL connection string
- `StaticDir`/`StaticUrl` - Static file serving configuration
- `TikaHost` - Apache Tika document parsing service
- `[Session]` - JWT secret key and session timeout
- `[Redis]` - Redis connection settings
- `[SMS]` - SMS provider configuration (Ali, Bao)
- `[OSS]` - Object storage configuration (Local, MinIO, Qiniu, Aliyun)
- `[ApiConfig]` - Third-party function API services

**Important:** Always change `Session.SecretKey` in production environments.

### Frontend Configuration (web/)

Environment variables:
- `.env.development` - Development environment config
- `.env.production` - Production environment config
- `VUE_APP_API_HOST` - Backend API server URL
- `VUE_APP_WS_HOST` - WebSocket server URL (auto-detected if empty)

## Adding New Features

### Backend Feature Development

1. **Define Data Model**: Create model struct in `api/store/model/`
2. **Create Value Object**: Define VO in `api/store/vo/` for API responses
3. **Implement Service**: Add business logic in `api/service/`
4. **Create Handler**: Add HTTP handler in `api/handler/` or `api/handler/admin/`
5. **Register Dependencies**: Wire up in `api/main.go` using `fx.Provide()` and `fx.Invoke()`
6. **Add Routes**: Routes are auto-registered in handlers via fx lifecycle

### Frontend Feature Development

1. **Create Component**: Add Vue component in `src/views/` or `src/components/`
2. **Add Route**: Register route in `src/router/index.js`
3. **Add Store (if needed)**: Create Pinia store in `src/store/`
4. **API Integration**: Add API calls in component or create service in `src/api/`
5. **Styling**: Add styles to module-specific `.styl` file in `src/assets/css/`, or to `common.styl` for reusable utilities

**CSS Priority**: Always check existing Stylus files before adding new styles. Prefer module-specific `.styl` files over inline component styles.

## Deployment

### Docker Deployment (Recommended)

1. Configure `deploy/conf/config.toml` with your settings
2. Configure Nginx in `deploy/conf/nginx/conf.d/`
3. Import initial database schema from `database/geekai_plus-v4.1.8.sql`
4. Run `cd deploy && docker-compose up -d`

**Services Exposed:**
- Port 80/443 - Web frontend (Nginx)
- Port 5678 - Backend API
- Port 3307 - MySQL
- Port 6380 - Redis

### Manual Deployment

1. Set up MySQL and Redis
2. Import database schema
3. Build backend: `cd api && make amd64`
4. Build frontend: `cd web && npm run build`
5. Configure and run backend: `cd api && ./geekai`
6. Serve frontend build from web server (Nginx recommended)

## Commit Message Guidelines

Follow these commit types as defined in the project:

- `feat:` - New features or functionality
- `fix:` - Bug fixes
- `docs:` - Documentation updates
- `style:` - Code style or component styling updates
- `refactor:` - Code refactoring without new features or bug fixes
- `opt:` - Performance optimization
- `chore:` - Small changes like text updates or comment modifications

## Important Notes

### Security Considerations
- Change JWT secret key (`Session.SecretKey`) in production
- Never commit sensitive credentials (API keys, database passwords)
- Payment gateway credentials require secure configuration
- Use environment variables for sensitive data

### External Service Dependencies
- MidJourney/Stable Diffusion require API access
- Payment services (Alipay, WeChat) need merchant accounts
- SMS services require provider credentials (Aliyun, etc.)
- OSS services require cloud storage accounts (Qiniu, Aliyun, MinIO)

### License and Compliance
- Apache 2.0 licensed
- Commercial use allowed but must retain copyright information
- Compliance with Chinese AI service regulations required for public deployment
- Do not provide generative AI services to Chinese public without proper filing (备案)
