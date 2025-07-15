
# DDD CRUD Server

A complete CRUD backend server built with Domain-Driven Design (DDD) architecture using Go Fiber and PostgreSQL.

## Architecture

This project follows DDD principles with clear separation of concerns:

- **Domain Layer**: Pure business logic (entities, repository interfaces)
- **Application Layer**: Use cases and application services
- **Infrastructure Layer**: External concerns (database, web, handlers)

## Project Structure

```
ddd-crud-server/
├── cmd/server/                 # Application entry point
├── configs/                    # Configuration management
├── internal/
│   ├── domain/                 # Business logic
│   ├── application/            # Use cases
│   └── infrastructure/         # External concerns
├── pkg/                        # Shared utilities
├── Dockerfile                  # Docker configuration
├── Makefile                    # Build automation
└── README.md
```

## Features

- ✅ Complete CRUD operations for users
- ✅ Clean DDD architecture
- ✅ Input validation
- ✅ Error handling
- ✅ Pagination support
- ✅ Structured logging
- ✅ Docker support
- ✅ Graceful shutdown

## API Endpoints

### Users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - Get all users (with pagination)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id`