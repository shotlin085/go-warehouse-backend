# Go Warehouse Backend - Simple School Project

A basic Go backend API for warehouse management with MySQL database.

## Features

- ✅ Health check endpoint with database status
- ✅ Create items (POST /api/items)
- ✅ Get all items (GET /api/items)
- ✅ Get single item (GET /api/items/:id)
- ✅ MySQL database in Docker container

## Setup

### 1. Start MySQL Docker Container

```bash
docker run -d \
  --name mysql-warehouse \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=warehouse_db \
  -p 3306:3306 \
  mysql:8.0
```

### 2. Install Go Dependencies

```bash
cd go-warehouse-backend
go mod download
```

### 3. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:7000`

## API Endpoints

### Health Check
```bash
GET http://localhost:7000/health
```

**Response:**
```json
{
  "status": "ok",
  "service": "go-warehouse-api",
  "database": {
    "type": "mysql",
    "connected": true,
    "message": "MySQL connected"
  },
  "time": "2024-01-15T10:30:00Z"
}
```

### Create Item
```bash
POST http://localhost:7000/api/items
Content-Type: application/json

{
  "name": "Widget A",
  "sku": "WDG-001",
  "description": "Test widget",
  "quantity": 10,
  "unit_price": 25.99
}
```

### Get All Items
```bash
GET http://localhost:7000/api/items
```

### Get Single Item
```bash
GET http://localhost:7000/api/items/1
```

## Environment Variables

Edit `.env` file:

```env
PORT=7000
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_DATABASE=warehouse_db
```

## Testing

### 1. Check Health
```bash
curl http://localhost:7000/health
```

### 2. Create an Item
```bash
curl -X POST http://localhost:7000/api/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Item",
    "sku": "TEST-001",
    "description": "My first item",
    "quantity": 5,
    "unit_price": 10.50
  }'
```

### 3. Get All Items
```bash
curl http://localhost:7000/api/items
```

## Docker Commands

**Start MySQL container:**
```bash
docker start mysql-warehouse
```

**Stop MySQL container:**
```bash
docker stop mysql-warehouse
```

**Check MySQL container status:**
```bash
docker ps | grep mysql-warehouse
```

**Remove MySQL container:**
```bash
docker rm -f mysql-warehouse
```

## Project Structure

```
go-warehouse-backend/
├── main.go          # Main application file
├── go.mod           # Go dependencies
├── .env             # Environment configuration
└── README.md        # This file
```

## Simple and Basic! 🎓

This is a simple school project with:
- Basic Go backend
- MySQL in Docker
- Health check showing database status
- Simple CRUD operations for items
