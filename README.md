# Subdub Backend

**Subdub** is a subscription tracking backend service, built in Go, that helps users monitor and manage their recurring subscriptions. It notifies users when subscriptions are about to expire so they can cancel before auto-renewal, preventing unwanted charges.

## 🚀 Features

- **User Authentication**
  - Sign up / Sign in with **username & password**
  - **JWT**-based authentication (Access & Refresh tokens)
  - **OAuth2** support (e.g., Google, GitHub)

- **CRUD Operations**
  - **Users**: Create, read, update, delete user profiles
  - **Subscriptions**: Register, view, update, and remove subscriptions

- **Automated Expiry Checks**
  - **Cron-style job** implemented using Go **goroutines**
  - Daily scans for subscriptions nearing expiry
  - Sends **email notifications** to users

- **Extensible & Secure**
  - Layered architecture for easy feature additions
  - Secure password hashing and token handling
  - Environment-based configuration

## 🏗️ Tech Stack

- **Language:** Go
- **Web Framework:** net/http (or Gin / Echo / Fiber)
- **Authentication:** JWT (github.com/dgrijalva/jwt-go) and OAuth2 (golang.org/x/oauth2)
- **Database:** PostgreSQL (via database/sql + pq) or your choice
- **Email:** SMTP via Go’s `net/smtp` or third-party service
- **Scheduling:** Go **goroutines** + `time.Ticker` for daily jobs

## 📐 Architecture Overview

```
+----------------+        +--------------+        +------------+
|  REST Clients  | <----> | API Endpoints| <----> | PostgreSQL |
+----------------+        +--------------+        +------------+
                               |    ^
                    JWT Auth   |    |  Cron Jobs (goroutines)
                               v    |
                           +--------+     Email Service
                           |  Auth  |----> (SMTP / SendGrid)
                           +--------+
```

## ⚙️ Getting Started

### Prerequisites

- Go 1.21+ installed
- PostgreSQL database
- SMTP credentials (for email notifications)

### Environment Variables

Create a `.env` file in the project root:

```dotenv
# Server
PORT=8080

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=password
DB_NAME=subdub

# JWT
JWT_ACCESS_SECRET=your_access_secret
JWT_REFRESH_SECRET=your_refresh_secret
JWT_ACCESS_EXP=15m
JWT_REFRESH_EXP=7d

# OAuth2 (e.g., Google)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# Email
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your_email@example.com
SMTP_PASS=email_password
```

### Installation & Run

```bash
# Clone the repo
git clone https://github.com/yourusername/subdub-backend.git
cd subdub-backend

# Install dependencies
go mod tidy

# Build
go build -o subdub-server ./cmd/server

# Run
./subdub-server
```

The server will be running at `http://localhost:8080`.

## 📄 API Endpoints

| Method | Endpoint                    | Description                                 |
|-------:|-----------------------------|---------------------------------------------|
| POST   | `/auth/signup`              | Register a new user                         |
| POST   | `/auth/login`               | Login with username/password                |
| GET    | `/auth/refresh`             | Refresh JWT tokens                          |
| GET    | `/auth/oauth/{provider}`    | Redirect to OAuth provider                  |
| GET    | `/auth/oauth/{provider}/cb` | OAuth callback                              |
| GET    | `/users/{id}`               | Retrieve user profile                       |
| PUT    | `/users/{id}`               | Update user info                            |
| DELETE | `/users/{id}`               | Delete user                                 |
| GET    | `/subs`                     | List all subscriptions for logged-in user   |
| POST   | `/subs`                     | Create a new subscription                   |
| GET    | `/subs/{id}`                | Get subscription details                    |
| PUT    | `/subs/{id}`                | Update subscription                         |
| DELETE | `/subs/{id}`                | Delete subscription                         |

## 🔄 Scheduled Tasks

A daily job runs at midnight via a Go **goroutine**:

- Scans subscriptions expiring in the next 3 days
- Sends reminder emails to users

Configured in `internal/cron` using `time.Ticker`.

## 🛡️ Security

- **Password Hashing:** bcrypt
- **Token Storage:** In-memory or persistent store (Redis) for refresh tokens
- **Rate Limiting:** Optional middleware

## 📦 Deployment

1. **Docker**
   ```dockerfile
   FROM golang:1.21 AS builder
   WORKDIR /app
   COPY . .
   RUN go mod tidy && go build -o subdub-server ./cmd/server

   FROM gcr.io/distroless/base
   COPY --from=builder /app/subdub-server /subdub-server
   CMD ["/subdub-server"]
   ```

2. **Kubernetes**
   - Create a Deployment and Service
   - Use ConfigMap & Secret for environment variables

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📜 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

*Happy Tracking!* 🚀

