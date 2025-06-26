# 📚 RESTful Library API

A production-ready REST API for managing a library's book collection, built with Go, MariaDB, and Docker.

## 🌟 Features

- **Complete CRUD Operations**: Create, read, update, and delete books
- **Advanced Search**: Search books by title or author
- **Pagination**: Efficient data retrieval with customizable page sizes
- **Health Monitoring**: Built-in health check endpoint
- **Structured Logging**: JSON-formatted logs with configurable levels
- **Containerized**: Full Docker and Docker Compose support
- **Production Ready**: Graceful shutdown, connection pooling, and error handling
- **CORS Support**: Cross-origin resource sharing for web clients

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

### Running with Docker Compose (Recommended)

1. **Clone and navigate to the project**:
   ```bash
   git clone <repository-url>
   cd library-api
   ```

2. **Start the application**:
   ```bash
   docker-compose up -d
   ```

3. **Verify the application is running**:
   ```bash
   curl http://localhost:8080/health
   ```

The API will be available at `http://localhost:8080` and the database at `localhost:3306`.

### Using Existing MariaDB Container

If you already have a MariaDB container running with the specified configuration:

```bash
docker run -d \
  --name=library-mariadb \
  -e MARIADB_DATABASE="db" \
  -e MARIADB_ROOT_PASSWORD="Password" \
  -e MARIADB_USER="user" \
  -e MARIADB_PASSWORD="Password" \
  -e ALLOW_EMPTY_PASSWORD=yes \
  -p 3306:3306  \
  -v "test-mariadb-vol:/var/lib/mysql"  \
  -d bitnami/mariadb-galera:10.11.4-debian-11-r0
```

Then run only the application:
```bash
docker build -t library-api .
docker run -p 8080:8080 \
  -e DB_HOST=localhost \
  -e DB_NAME=db \
  -e DB_USER=user \
  -e DB_PASSWORD=Password \
  --network host \
  library-api
```

### Local Development Setup

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your database settings
   ```

3. **Ensure your MariaDB container is running** with the configuration above.

4. **Run the application**:
   ```bash
   go run main.go
   ```

## 📖 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### Health Check
```http
GET /health
```
Returns application health status.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### List Books
```http
GET /api/v1/books?page=1&limit=10&q=search_term
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page, max 100 (default: 10)
- `q` (optional): Search term for title or author

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "The Go Programming Language",
      "author": "Alan Donovan, Brian Kernighan",
      "published_year": 2015,
      "available": true,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

#### Get Single Book
```http
GET /api/v1/books/{id}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "The Go Programming Language",
    "author": "Alan Donovan, Brian Kernighan",
    "published_year": 2015,
    "available": true,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

#### Create Book
```http
POST /api/v1/books
Content-Type: application/json

{
  "title": "New Book Title",
  "author": "Author Name",
  "published_year": 2024,
  "available": true
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 11,
    "title": "New Book Title",
    "author": "Author Name",
    "published_year": 2024,
    "available": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "message": "Book created successfully"
}
```

#### Update Book
```http
PUT /api/v1/books/{id}
Content-Type: application/json

{
  "title": "Updated Title",
  "available": false
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Updated Title",
    "author": "Alan Donovan, Brian Kernighan",
    "published_year": 2015,
    "available": false,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:35:00Z"
  },
  "message": "Book updated successfully"
}
```

#### Delete Book
```http
DELETE /api/v1/books/{id}
```

**Response:**
```json
{
  "success": true,
  "message": "Book deleted successfully"
}
```

### Error Responses

All error responses follow this format:
```json
{
  "success": false,
  "error": "Error description"
}
```

Common HTTP status codes:
- `400` - Bad Request (invalid input)
- `404` - Not Found (book doesn't exist)
- `500` - Internal Server Error

## 🏗️ Architecture

The application follows clean architecture principles:

```
├── main.go              # Application entry point
├── handlers/            # HTTP handlers
│   └── book_handler.go  # Book-related endpoints
├── models/              # Data models and DTOs
│   └── book.go          # Book model and request/response types
├── db/                  # Database layer
│   └── db.go            # Database operations and migrations
├── migrations/          # Database schema files
│   └── 01_init.sql      # Initial schema and sample data
├── Dockerfile           # Container configuration
├── docker-compose.yml   # Multi-container orchestration
└── README.md           # Documentation
```

## 🛠️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `3306` |
| `DB_NAME` | Database name | `db` |
| `DB_USER` | Database user | `user` |
| `DB_PASSWORD` | Database password | `Password` |
| `PORT` | Application port | `8080` |
| `LOG_LEVEL` | Log level (debug, info, warn, error) | `info` |

### Database Schema

The application automatically creates the required database schema on startup. The `books` table includes:

- Optimized indexes for query performance
- Automatic timestamps for audit trails
- Data validation constraints
- Sample data for testing

## 🔧 Development

### Building from Source

```bash
# Build for current platform
go build -o library-api

# Build for Linux (for Docker)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Docker Commands

```bash
# Build image
docker build -t library-api .

# Run container with existing MariaDB
docker run -p 8080:8080 \
  -e DB_HOST=localhost \
  -e DB_NAME=db \
  -e DB_USER=user \
  -e DB_PASSWORD=Password \
  --network host \
  library-api

# View logs
docker logs library-api

# Stop services
docker-compose down
```

## 📊 Performance Features

- **Connection Pooling**: Configured for optimal database performance
- **Pagination**: Efficient handling of large datasets
- **Indexing**: Strategic database indexes for fast queries
- **Graceful Shutdown**: Proper cleanup of resources
- **Health Checks**: Built-in monitoring for container orchestration

## 🔒 Security Features

- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries
- **CORS Configuration**: Configurable cross-origin policies
- **Non-root Container**: Security-hardened Docker image

## 🚢 Production Deployment

For production deployment:

1. **Use environment-specific configurations**
2. **Set up proper logging aggregation**
3. **Configure database backups**
4. **Implement monitoring and alerting**
5. **Use secrets management for sensitive data**
6. **Set up load balancing for high availability**

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the API documentation above
- Review the logs using `docker logs library-api`

## 🎯 Future Enhancements

- [ ] Authentication and authorization
- [ ] Rate limiting
- [ ] Caching layer (Redis)
- [ ] Full-text search (Elasticsearch)
- [ ] API versioning
- [ ] OpenAPI/Swagger documentation
- [ ] Metrics and monitoring (Prometheus)
- [ ] Integration tests
- [ ] CI/CD pipeline