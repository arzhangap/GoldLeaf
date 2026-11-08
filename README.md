# 🍃 GoldLeaf

A lightweight paste-sharing web application built with Go.

GoldLeaf provides a simple way to create, store, and share text snippets through unique URLs, with support for expiration dates.

## ✨ Features

- **Create pastes** — Save text snippets with an optional title.
- **Unique URLs** — Access saved pastes through generated codes.
- **Expiration support** — Set an expiration period for your pastes.
- **Web interface** — Browse and view pastes through a clean interface.
- **Persistent storage** — Store paste data in a MySQL database.

## 🛠️ Tech Stack

| Technology | Purpose |
|---|---|
| Go | Backend |
| MySQL | Database |
| httprouter | HTTP routing |
| HTML / CSS | Web interface |

## 🚀 Getting Started

### Prerequisites

- Go 1.26.5 or later
- MySQL

### Installation

Clone the repository:

```bash
git clone https://github.com/arzhangap/GoldLeaf.git
cd GoldLeaf
```

Download dependencies:

```bash
go mod download
```

Configure your MySQL database connection according to the application's configuration.

Run the application using its configured entry point.

## 📁 Project

[GoldLeaf on GitHub](https://github.com/arzhangap/GoldLeaf)

## 📄 License

See the repository for licensing information.
