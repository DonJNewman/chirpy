# Chirpy - Social Media API # 

This professional grade backend API is built with Go and PostgreSQL and features user management, post creation, and comprehensive metrics tracking.

This project is still in progress. Future functionality includes authentication and authorization. 

## 🚀 Key Features

- **RESTful API**: Clean HTTP endpoints for social media operations
- **User Management**: Secure user registration and authentication
- **Post System**: Create, retrieve, and manage user posts (chirps)
- **Metrics Tracking**: Thread-safe request counting and performance monitoring
- **Database Integration**: PostgreSQL with type-safe query generation

## 🎯 Skills Demonstrated

- **Web API Development**: RESTful endpoint design and HTTP protocol mastery
- **Database Integration**: PostgreSQL schema design and query optimization
- **Concurrent Programming**: Thread-safe operations and atomic variables
- **Security Practices**: Environment-based configuration and input validation
- **Code Generation**: SQLC integration for type-safe database operations
- **Middleware Architecture**: Custom middleware for cross-cutting concerns
- **Error Handling**: Comprehensive error responses and status codes


## 🏗️ Technical Architecture

### Web Server Design
- **Standard Library HTTP**: Custom middleware and routing without external frameworks
- **Middleware Pipeline**: Authentication, CORS, and metrics collection
- **Environment Configuration**: Secure configuration management with godotenv
- **Static File Serving**: Integrated file server for frontend assets

### Database Architecture
- **PostgreSQL Backend**: Relational database with proper schema design
- **SQLC Integration**: Compile-time SQL query validation and type generation
- **Migration System**: Schema versioning and database evolution
- **Connection Management**: Efficient database connection handling

### Concurrent Safety
- **Atomic Operations**: Thread-safe metrics counting using sync/atomic
- **Goroutine-Safe Design**: Concurrent request handling with proper synchronization
- **Resource Management**: Proper cleanup and connection pooling

## 💻 Technologies Used

- **Go 1.21+**: Modern Go with generics and latest features
- **PostgreSQL**: Enterprise-grade relational database
- **SQLC**: Type-safe SQL query generation
- **Godotenv**: Environment variable management
- **Google UUID**: Secure unique identifier generation
- **Standard Library HTTP**: Zero-dependency web server implementation

## 🔧 Technical Highlights

### Custom HTTP Middleware
```go
// Authentication middleware with proper error handling
func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Token validation and user context injection
    }
}
```

### Thread-Safe Metrics
```go
// Atomic operations for concurrent metric tracking
atomic.AddInt64(&cfg.fileserverHits, 1)
```

### Type-Safe Database Queries
```go
// SQLC-generated queries with compile-time validation
user, err := cfg.db.CreateUser(ctx, database.CreateUserParams{
    ID:        userID,
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
    Email:     email,
})
```

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL 12+

### Environment Setup
Create a `.env` file:
```bash
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
```

### Installation & Running
```bash
# Install dependencies
go mod download

# Run database migrations
# (migrations handled automatically)

# Start the server
go run .
# Server runs on http://localhost:8080
```

## 📋 API Endpoints

### Health & Metrics
- `GET /api/healthz` - Health check endpoint
- `GET /admin/metrics` - Request metrics and server statistics
- `POST /admin/reset` - Reset metrics (dev mode only)

### User Management
- `POST /api/users` - Create new user account
- `PUT /api/users` - Update user information (authenticated)

### Posts (Chirps)
- `POST /api/chirps` - Create new post (authenticated)
- `GET /api/chirps` - Retrieve all posts with optional filtering
- `GET /api/chirps/{chirpID}` - Get specific post by ID
- `DELETE /api/chirps/{chirpID}` - Delete post (owner only)

### Static Files
- `GET /` - Serve static web assets
- File server with request tracking and metrics

## 🏁 Project Structure

```
chirpy/
├── internal/database/     # SQLC-generated database code
├── sql/                   # SQL migrations and queries
├── assets/               # Static web assets
├── *.go                 # HTTP handlers and middleware
└── .env                 # Environment configuration
```

This project demonstrates professional-grade web API development with modern Go practices, suitable for roles requiring scalable web service development.
