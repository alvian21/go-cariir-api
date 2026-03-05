# Go Cariir API

A robust and scalable Restful API built with Golang using the Fiber framework, GORM for ORM, and PostgreSQL.

## 🚀 Features

- **Authentication**: Secure JWT-based authentication.
- **RBAC**: Role-Based Access Control with support for granular permissions.
- **Permission Guard**: Middleware for protecting routes based on specific permissions.
- **Database Migrations**: Managed with Goose.
- **Seeding**: Automated database seeding for initial data.
- **Standardized Responses**: Consistent JSON structure for all API responses.
- **Live Reload**: Development support using Air.

## 🛠 Tech Stack

- **Language**: [Golang](https://golang.org/)
- **Framework**: [Fiber](https://gofiber.io/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Migration Tool**: [Goose](https://github.com/pressly/goose)
- **Validation**: [validator/v10](https://github.com/go-playground/validator)

## 📦 Project Structure

```text
.
├── cmd/                # Entry points (Main, Reseters, Migrations)
├── config/             # Configuration management
├── database/           # Database connections, migrations, and seeders
├── handler/            # HTTP request handlers
├── middleware/         # Custom middlewares (Auth, etc.)
├── model/              # Entities, Requests, and Responses models
├── route/              # API route definitions
├── utils/              # Helper functions and utilities
└── main.go             # Application entry point
```

## ⚙️ Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- [Air](https://github.com/cosmtrek/air) (for development)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/alvian21/go-cariir-api.git
   cd go-cariir-api
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Setup environment variables:

   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. Run migrations and seed the database:

   ```bash
   go run main.go --migrate --seed
   ```

   _(Alternatively, use the specific migration/reset commands provided in the `cmd/` directory)_

5. Run the application:

   ```bash
   # Using Air (recommended for development)
   air

   # Standard Go run
   go run main.go
   ```

## 🛡 RBAC & Middleware

This project implements a granular Role-Based Access Control (RBAC) system.

### Auth Middleware

Protects routes and ensures the user is authenticated. It stores user claims in `context.Locals("user")`.

### Permission Guard

Protects routes based on specific permissions.

```go
user.Post("/", middleware.PermissionGuard("user.create"), handler.UserHandlerCreate)
```

## 🛣 API Endpoints

### Auth

- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### User

- `GET /api/users` - List all users (Requires Auth)
- `GET /api/users/:id` - Get user details

## 📝 License

This project is open-source and available under the MIT License.
