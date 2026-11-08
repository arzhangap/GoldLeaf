# 🍃 GoldLeaf

**A lightweight paste-sharing application built with Go and MySQL.**

GoldLeaf lets you create text snippets and share them through unique, easy-to-use URLs. It provides a simple web interface for creating and retrieving pastes, with configurable expiration periods.

## Features

- **Create pastes** — Save text snippets with an optional title.
- **Unique URLs** — Each paste receives a randomly generated six-character code.
- **Automatic expiration** — Choose an expiration period of 1 hour, 1 day, or 7 days.
- **Persistent storage** — Store paste data in MySQL.
- **Server-rendered HTML** — Render pages using Go's `html/template` package.
- **Static asset embedding** — Serve application assets through Go's embedded filesystem.

## Screenshots

<!-- Add screenshots to docs/screenshots/ and uncomment these lines.

| Home | Create a paste |
|:---:|:---:|
| ![Home page](docs/screenshots/home.png) | ![Create paste page](docs/screenshots/create-paste.png) |

![Paste view](docs/screenshots/read-paste.png)

-->

## Tech Stack

| Technology | Purpose |
|---|---|
| Go | Application backend |
| MySQL | Persistent data storage |
| `database/sql` | Database access |
| `go-sql-driver/mysql` | MySQL driver |
| `httprouter` | HTTP routing |
| `html/template` | Server-side HTML rendering |
| HTML, CSS, JavaScript | Frontend |

## Project Structure

```text
GoldLeaf/
├── cmd/
│   └── web/              # HTTP server, routes and handlers
├── internal/
│   ├── models/           # Database models and queries
│   └── validator/        # Form validation
├── ui/
│   ├── html/             # HTML templates
│   └── static/           # CSS and JavaScript
├── tls/                  # TLS certificate and private key
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Go 1.26.5 or later
- MySQL
- OpenSSL (for generating a local TLS certificate)

### 1. Clone the repository

```bash
git clone https://github.com/arzhangap/GoldLeaf.git
cd GoldLeaf
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure MySQL

Create a database:

```sql
CREATE DATABASE goldleaf
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;
```

Create the `pastes` table:

```sql
USE goldleaf;

CREATE TABLE pastes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(6) NOT NULL UNIQUE,
    title VARCHAR(100) NOT NULL,
    body TEXT NOT NULL,
    password_hash VARCHAR(255) NULL,
    created DATETIME NOT NULL,
    expires DATETIME NULL
);
```

The application expects a MySQL user with permission to access this database.

### 4. Configure TLS

GoldLeaf currently starts its HTTP server using TLS and expects the following files:

```text
tls/
├── cert.pem
└── key.pem
```

For local development, generate a self-signed certificate:

```bash
mkdir -p tls

openssl req -x509 -newkey rsa:2048 \
    -keyout tls/key.pem \
    -out tls/cert.pem \
    -days 365 \
    -nodes \
    -subj "/CN=localhost"
```

This certificate is intended for local development, not production deployment.

### 5. Run the application

```bash
go run ./cmd/web \
    -addr=":4000" \
    -dsn="USER:PASSWORD@tcp(127.0.0.1:3306)/goldleaf?parseTime=true"
```

Replace `USER` and `PASSWORD` with your MySQL credentials.

Open:

```text
https://localhost:4000
```

Your browser will warn you about the self-signed certificate. This is expected during local development.

## HTTP Routes

| Method | Route | Description |
|---|---|---|
| `GET` | `/` | Home page |
| `GET` | `/paste/create` | Display the paste creation form |
| `POST` | `/paste/create` | Create a new paste |
| `GET` | `/paste/view/:code` | Display an existing paste |
| `POST` | `/paste/view` | Submit a paste lookup |
| `GET` | `/static/*filepath` | Serve static assets |

## Configuration

The application accepts command-line flags:

| Flag | Default | Description |
|---|---|---|
| `-addr` | `:4000` | HTTP server address |
| `-dsn` | `goland:pass@/goldleaf?parseTime=true` | MySQL data source name |

For deployments, provide database credentials through your deployment configuration rather than committing them to source control.

## 📄 License

Copyright (c) 2026 arzhangap

This project is licensed under the MIT License.
See the [LICENSE](LICENSE) file for details.
