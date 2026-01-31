# Project 27

A robust backend API for a note-taking application, built with Go (Golang) and PostgreSQL. This project features secure user authentication, folder organization, and note management using a clean architecture approach.

## Features

- **User Management**: Registration, profile retrieval, updates, and deletion.
- **Authentication**: Secure JWT-based authentication with bcrypt password hashing.
- **Folder Organization**: Create, list, and delete folders to organize content.
- **Note Management**: Retrieve notes associated with specific folders.
- **Transactional Integrity**: Critical database operations use transactions to ensure data consistency.

## Tech Stack

- **Language**: Go
- **Router**: [chi](https://github.com/go-chi/chi) - lightweight, idiomatic and composable router.
- **Database**: PostgreSQL
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Logging**: [zap](https://github.com/uber-go/zap) - Blazing fast, structured logging.

## Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL instance
- `migrate` CLI tool installed

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/MihirSahani/Project-27.git
   cd Project-27
   ```

2. **Environment Configuration**
   Set the `POSTGRES_ADDRESS` environment variable required for the application and migrations.
   ```bash
   export POSTGRES_ADDRESS="postgres://user:password@localhost:5432/dbname?sslmode=disable"
   ```

3. **Database Migrations**
   Use the included helper script to run migrations.
   ```bash
   chmod +x cmd.sh
   ./cmd.sh up
   ```

4. **Run the Application**
   ```bash
   go run ./api
   ```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| **Auth** | | |
| POST | `/authenticate` | Login and receive a JWT token |
| **Users** | | |
| POST | `/user` | Register a new user |
| GET | `/user/{id}` | Get user details |
| PUT | `/user/{id}` | Update user profile |
| DELETE | `/user/{id}` | Delete a user account |
| **Folders** | | |
| POST | `/folder` | Create a new folder |
| GET | `/folder` | Get all folders for the current user |
| DELETE | `/folder/{id}` | Delete a specific folder |
| GET | `/folder/{id}/note` | Get all notes inside a specific folder |

## Project Structure

- **`api/`**: Contains HTTP handlers, middleware, and routing logic.
- **`storage/`**: Defines the storage interfaces and entity structs.
  - **`postgres/`**: Concrete implementation of storage interfaces using PostgreSQL.
- **`internal/`**: Shared utilities, validation logic, and error definitions.
- **`cmd.sh`**: Shell script wrapper for managing database migrations.