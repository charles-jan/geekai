# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GeekAI is a ChatGPT Plus backend API implemented in Go, featuring a Gin + MySQL architecture with dependency injection using the fx framework and GORM ORM. This is a comprehensive AI platform supporting multiple services including ChatGPT, MidJourney, Stable Diffusion, DALL-E, Suno music generation, and video generation.

## Build and Development Commands

### Build Commands
```bash
# Build for Linux AMD64
make amd64

# Build for Linux ARM64  
make arm64

# Build both architectures
make

# Clean build artifacts
make clean
```

### Development Commands
```bash
# Run the application
go run main.go

# Run tests
go test ./test/...

# Run a specific test
go test ./test/service/captcha_service_test.go -v

# Generate test coverage
go test -cover ./...
```

### Configuration
- Copy `config.sample.toml` to `config.toml` and configure your settings
- Set `CONFIG_FILE` environment variable to use a different config file
- Set `APP_DEBUG=true` for debug mode

## Architecture Overview

### Core Architecture
The application uses **uber-go/fx** for dependency injection and follows a layered architecture:

- **main.go**: Application entry point with fx dependency injection setup
- **core/**: Core application server, middleware, and configuration
- **handler/**: HTTP request handlers (controllers)
- **service/**: Business logic services
- **store/**: Data persistence layer (models, database connections)
- **utils/**: Utility functions and helpers

### Key Components

#### Dependency Injection (FX Framework)
The entire application is wired using Uber's fx framework. All services, handlers, and providers are registered in `main.go` using `fx.Provide()` and initialized with `fx.Invoke()`.

#### Database Layer
- **MySQL**: Primary database using GORM
- **Redis**: Caching and session storage
- **LevelDB**: Embedded key-value storage

#### Service Architecture
Each major feature is organized as a service:
- **dalle/**: DALL-E image generation service
- **mj/**: MidJourney image generation service  
- **sd/**: Stable Diffusion service
- **suno/**: Music generation service
- **video/**: Video generation service (Luma)
- **payment/**: Payment processing (Alipay, WeChat, etc.)
- **oss/**: Object storage services (Local, MinIO, Qiniu, Aliyun)
- **sms/**: SMS services (Aliyun, etc.)

#### Handler Organization
- **handler/**: Public API endpoints
- **handler/admin/**: Admin panel API endpoints

### Configuration System
- TOML-based configuration in `config.toml`
- Environment variable overrides supported
- Database config cached in `types.SystemConfig`

### Authentication & Authorization
- JWT token-based authentication
- Separate admin and user authentication flows
- Redis-based session management
- Middleware handles authorization in `core/app_server.go:authorizeMiddleware()`

### WebSocket Support
- Real-time communication support
- Chat streaming functionality
- Token-based WebSocket authentication

## Directory Structure

```
/
├── core/           # Application server and core functionality
│   ├── app_server.go    # Main server with middleware
│   ├── config.go        # Configuration loading
│   └── types/           # Core type definitions
├── handler/        # HTTP request handlers
│   ├── admin/           # Admin panel handlers
│   └── *.go            # Public API handlers
├── service/        # Business logic services
│   ├── dalle/          # DALL-E service
│   ├── mj/             # MidJourney service
│   ├── payment/        # Payment services
│   └── *.go           # Other services
├── store/          # Data persistence layer
│   ├── model/          # Database models
│   ├── vo/             # Value objects
│   ├── mysql.go        # MySQL connection
│   └── redis.go        # Redis connection
├── utils/          # Utility functions
│   └── resp/          # HTTP response helpers
├── test/           # Test files
├── res/            # Static resources
├── logs/           # Application logs
└── main.go         # Application entry point
```

## Development Guidelines

### Adding New Features
1. Define models in `store/model/`
2. Create value objects in `store/vo/`
3. Implement service logic in `service/`
4. Create handlers in `handler/` or `handler/admin/`
5. Register dependencies in `main.go` using fx.Provide/fx.Invoke

### Database Operations
- Use GORM for database operations
- Follow the existing model pattern in `store/model/`
- Use transactions for complex operations

### API Development
- Follow RESTful conventions
- Use the response helpers in `utils/resp/`
- Implement proper error handling
- Add authentication where required

### Testing
- Write unit tests following the pattern in `test/service/captcha_service_test.go`
- Use dependency injection for testable code
- Mock external dependencies

## Important Notes

### Security Considerations
- JWT secret keys should be changed in production
- Database credentials are in config files
- Payment gateway credentials require secure configuration

### Service Dependencies
Many services have external dependencies:
- MidJourney requires API access
- Payment services need merchant accounts
- SMS services need provider credentials
- OSS services need cloud storage accounts

### Background Jobs
- Uses XXL-Job for scheduled tasks
- Services run background goroutines for processing
- Check service Run() methods for background operations