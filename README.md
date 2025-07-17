# Airbnb Clone


## 🏗️ Architecture Overview

This system follows clean architecture principles with:
- **Repository Layer**: Data access abstraction with both GORM and raw SQL
- **Service Layer**: Business logic implementation
- **API Layer**: RESTful HTTP handlers using Gin framework
- **Middleware**: Authentication, CORS, logging, and rate limiting
- **Database**: PostgreSQL with optimized indexes

## 🚀 Features

### Core Functionality
- **User Management**: Registration and authentication
- **Property Listings**: CRUD operations for property management
- **Booking System**: Complete booking workflow with availability checking
- **Review System**: Property reviews and ratings
- **Role-based Access**: Guest, Host, and Admin roles

## 🛠️ Technologies

- **Language**: Go 1.21
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens
- **Testing**: Go testing framework

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 13+ 
- Git

## 🚀 Quick Start

### 1. Clone the repository
```bash
git clone <repository-url>
cd airbnb-clone
```

### 2. Install dependencies
```bash
go mod download
```

### 3. Set up environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

### 4. Set up PostgreSQL database
```bash
# Create database
createdb airbnb_clone

# Or using PostgreSQL client
psql -c "CREATE DATABASE airbnb_clone;"
```

### 5. Run the application
```bash
go run cmd/server/main.go
```

## How to run on Docker

### 1. Clone the repository
```bash
git clone <repository-url>
cd airbnb-clone
```

### 2. Create Your .env File
```bash
PORT=8081
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=airbnb_clone
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=your-super-secret-jwt-key-for-development
ENVIRONMENT=development

```

### 3. 🐳 Start all services (PostgreSQL, Redis, and the Go app):
```bash
docker-compose up --build
```
