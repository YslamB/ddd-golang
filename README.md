
# DDD CRUD Server

A complete CRUD backend server built with Domain-Driven Design (DDD) architecture using Go Fiber and PostgreSQL.

## Architecture

This project follows DDD principles with clear separation of concerns:

- **Domain Layer**: Pure business logic (entities, repository interfaces)
- **Application Layer**: Use cases and application services
- **Infrastructure Layer**: External concerns (database, web, handlers)

## Project Structure

```
ddd-golang/
  ├── cmd/
  │   └── server/
  │       └── main.go
  ├── configs/
  │   └── config.go
  ├── internal/
  │   ├── application/
  │   │   ├── dto/
  │   │   │   ├── request_dto.go
  │   │   │   └── user_dto.go
  │   │   └── service/
  │   │       └── user_service.go
  │   ├── domain/
  │   │   ├── entity/
  │   │   │   └── user.go
  │   │   └── repository/
  │   │       └── user_repository.go
  │   └── infrastructure/
  │       ├── database/
  │       │   └── postgres.go
  │       ├── handler/
  │       │   └── user_handler.go
  │       ├── repository/
  │       │   └── user_repository.go
  │       └── web/
  │           ├── error_handler.go
  │           ├── middleware/
  │           │   ├── cors.go
  │           │   └── logger.go
  │           └── routes.go
  ├── pkg/
  │   ├── logger/
  │   │   └── logger.go
  │   ├── utils/
  │   │   └── response.go
  │   └── validator/
  │       └── validator.go
  ├── go.mod
  ├── go.sum
  ├── Makefile
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