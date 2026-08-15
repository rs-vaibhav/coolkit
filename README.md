<div align="center">

# 🚀 CoolKit

**The Operating System for Student Clubs**

[![CI](https://github.com/coolkit-org/coolkit/actions/workflows/ci.yml/badge.svg)](https://github.com/coolkit-org/coolkit/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

Centralize your club management. Events, members, finances, tasks — all in one modern platform.

[Getting Started](#-quickstart) · [API Reference](docs/api.md) · [Contributing](CONTRIBUTING.md) · [Architecture](docs/architecture.md)

</div>

---

## ✨ Features

- 👥 **Member Management** - Add, remove, and manage club members.
- 🏢 **Multi-Club Support** - Manage multiple clubs from a single account.
- 🔐 **Secure Authentication** - JWT-based auth with bcrypt password hashing.
- 🚀 **Blazing Fast** - Built with Go and Gin for maximum performance.
- 📊 **PostgreSQL Backed** - Reliable data storage with GORM.

## 🏗 Architecture

```text
Client Request
      │
      ▼
┌──────────────┐
│  Middleware  │ (Auth, Logger, CORS, Recovery)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Handlers   │ (HTTP request/response, validation)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Services   │ (Business logic, transactions)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Repositories │ (Data access layer, GORM queries)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  PostgreSQL  │ (Database)
└──────────────┘
```

## 🛠 Tech Stack

| Layer | Technology | Why |
|-------|------------|-----|
| Language | Go 1.22+ | Performance, concurrency, strong typing. |
| Web Framework | Gin | Fast routing, excellent middleware support. |
| Database | PostgreSQL | Relational integrity, JSONB support, reliability. |
| ORM | GORM | Developer productivity, migrations, associations. |
| Auth | JWT & bcrypt | Stateless authentication, secure password storage. |
| Config | Viper | Multi-format configuration management, env vars. |

## 🚀 Quickstart

### Option 1: Docker (Recommended)

1. Set up environment:
   ```bash
   cp .env.example .env
   ```
2. Start the application:
   ```bash
   make docker-up
   ```
3. Visit `http://localhost:8080/api/v1/health`

### Option 2: Manual Setup

**Prerequisites:** Go 1.22+, PostgreSQL 14+

1. Clone the repository:
   ```bash
   git clone https://github.com/coolkit-org/coolkit.git
   cd coolkit
   ```
2. Set up PostgreSQL:
   - Create a database (e.g., `coolkit`)
   - Ensure your PostgreSQL server is running
3. Configure environment:
   ```bash
   cp .env.example .env
   # Edit .env with your DB credentials
   ```
4. Run migrations and seed data:
   ```bash
   make migrate-up
   make seed
   ```
5. Start the server:
   ```bash
   make run
   ```

## 🔌 API Overview

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/health` | 🔓 | Basic health check |
| GET | `/api/v1/health/db` | 🔓 | Database health check |
| POST | `/api/v1/auth/register` | 🔓 | Register a new user |
| POST | `/api/v1/auth/login` | 🔓 | Authenticate and get JWT |
| GET | `/api/v1/auth/me` | 🔒 | Get current user profile |
| POST | `/api/v1/clubs` | 🔒 | Create a new club |
| GET | `/api/v1/clubs` | 🔒 | List user's clubs |
| GET | `/api/v1/clubs/:id` | 🔒 | Get club details |
| POST | `/api/v1/clubs/join` | 🔒 | Join a club using a join code |
| GET | `/api/v1/clubs/:id/members`| 🔒 | List members of a club |
| PUT | `/api/v1/clubs/:id/members/:user_id/role`| 🔒 | Update member role |
| DELETE | `/api/v1/clubs/:id/members/:user_id`| 🔒 | Remove a member |
| DELETE | `/api/v1/clubs/:id/members/me`| 🔒 | Leave a club |
| POST | `/api/v1/clubs/:id/announcements` | 🔒 | Post an announcement |
| GET | `/api/v1/clubs/:id/announcements` | 🔒 | List announcements |
| DELETE | `/api/v1/announcements/:id` | 🔒 | Delete an announcement |
| POST | `/api/v1/clubs/:id/events` | 🔒 | Create an event |
| GET | `/api/v1/clubs/:id/events` | 🔒 | List events |
| GET | `/api/v1/events/:id` | 🔒 | Get event details |
| PUT | `/api/v1/events/:id` | 🔒 | Update an event |
| DELETE | `/api/v1/events/:id` | 🔒 | Delete an event |
| POST | `/api/v1/events/:id/roles` | 🔒 | Assign an event role |
| GET | `/api/v1/events/:id/roles` | 🔒 | List event roles |
| DELETE | `/api/v1/events/:id/roles/:role_id`| 🔒 | Remove an event role |
| POST | `/api/v1/events/:id/tasks` | 🔒 | Create a task |
| GET | `/api/v1/events/:id/tasks` | 🔒 | List tasks |
| PATCH | `/api/v1/tasks/:id/status` | 🔒 | Update task status |
| DELETE | `/api/v1/tasks/:id` | 🔒 | Delete a task |
| POST | `/api/v1/events/:id/finance` | 🔒 | Log finance entry |
| GET | `/api/v1/events/:id/finance` | 🔒 | Get finance summary |
| DELETE | `/api/v1/finance/:id` | 🔒 | Delete finance entry |

*For full details, see the [API Reference](docs/api.md).*

## 📁 Project Structure

```text
.
├── cmd
│   └── server          # Application entrypoint
├── docs                # Documentation
├── internal
│   ├── config          # Configuration loading
│   ├── database        # DB connection and migrations
│   ├── handler         # HTTP handlers
│   ├── middleware      # Gin middlewares
│   ├── model           # Domain models
│   ├── repository      # Data access layer
│   └── service         # Business logic
├── pkg
│   └── response        # Standardized HTTP responses
└── scripts
    └── seed            # Database seeding script
```

## 💻 Development

See our [Development Guide](docs/development.md) for detailed instructions on local development, running tests, and database migrations.

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on how to get started.

Look out for issues labeled [`good first issue`](https://github.com/coolkit-org/coolkit/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) if you're new to the project!

## 🗺 Roadmap

### Phase 1: Foundation (Completed)
- [x] Project setup and architecture
- [x] User authentication (JWT)
- [x] Basic club management
- [x] API Documentation

### Phase 2: Engagement (Completed)
- [x] Event scheduling and management
- [x] RSVP system (via Join Codes and Join Requests)
- [x] Announcements and notifications
- [x] Member roles and permissions (Admin/Coordinator/Member)

### Phase 3: Operations (Completed)
- [x] Finance tracking (dues, expenses)
- [x] Task management
- [ ] Resource booking
- [ ] Analytics dashboard

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---
<div align="center">
Built with ❤️ by FOSS Club, SRM IST
</div>
