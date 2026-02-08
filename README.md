# Cinema System - GoLang Backend

Cinema booking and management system built with Go, Gin, and MongoDB.

## Team Members
- Mukhammedali Kalen
- Shyngys Abdullaev
- Alikhan Orynbasarov

## Features
- User authentication and authorization (User/Admin roles)
- Movie catalog with genres and reviews
- Session scheduling and hall management
- Ticket booking with seat selection
- Simple payment system (balance-based)
- Review and rating system with automatic calculation
- Concurrent booking handling with goroutines

## Tech Stack
- **Framework**: Gin
- **Database**: MongoDB
- **Language**: Go 1.21+

## Project Structure
```
/cinema-system
├── main.go
├── go.mod
├── internal/
│   ├── config/        # Database configuration
│   ├── models/        # Data models
│   ├── repositories/  # Database operations
│   ├── services/      # Business logic
│   ├── handlers/      # HTTP handlers
│   ├── middleware/    # Auth middleware
│   └── routes/        # Route definitions
```

## Installation

1. Install MongoDB:
```bash
# Ubuntu/Debian
sudo apt install mongodb

# macOS
brew install mongodb-community
```

2. Clone and setup:
```bash
git clone <repository-url>
cd cinema-system
go mod tidy
```

3. Configure environment:
```bash
cp .env.example .env
# Edit .env with your settings
```

4. Run the application:
```bash
go run main.go
```

## API Endpoints

### Public
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `GET /api/movies` - Get all movies
- `GET /api/movies/:id` - Get movie details
- `GET /api/sessions/upcoming` - Get upcoming sessions
- `GET /api/sessions/movie/:movieId` - Get sessions for movie

### User (Requires Auth)
- `POST /api/bookings` - Book ticket
- `DELETE /api/bookings/:id` - Cancel booking
- `GET /api/bookings/my` - Get my tickets
- `POST /api/reviews` - Create review
- `GET /api/reviews/movie/:movieId` - Get movie reviews

### Admin (Requires Admin Role)
- `POST /api/admin/movies` - Create movie
- `PUT /api/admin/movies/:id` - Update movie
- `DELETE /api/admin/movies/:id` - Delete movie
- `POST /api/admin/halls` - Create hall
- `POST /api/admin/sessions` - Create session
- `GET /api/admin/bookings` - View all bookings
- `DELETE /api/admin/reviews/:id` - Delete review

## Authentication

For testing, pass headers:
```
X-User-ID: <mongodb-object-id>
X-User-Role: USER or ADMIN
```

## Database Models

### User
- Roles: GUEST, USER, ADMIN
- Balance for payments
- Password hashing with bcrypt

### Movie
- Auto-calculated rating from reviews
- Genre associations

### Session
- Auto-calculated end time
- Hall and movie associations

### Ticket
- Statuses: BOOKED, PAID, CANCELLED
- Seat validation
- Balance deduction

## Concurrency Features

The system uses Go goroutines for:
- Async rating calculation after reviews
- Parallel ticket creation and payment
- Batch rating updates
- Concurrent booking validation

## Testing

Run with sample data:
```bash
# Start MongoDB
sudo systemctl start mongodb

# Run application
go run main.go

# Test registration
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone_number": "+77771234567",
    "password": "password123"
  }'
```

## Future Enhancements
- JWT authentication
- WebSocket for real-time seat updates
- Email notifications
- Payment gateway integration
- Advanced search and filtering
- Analytics dashboard

## License
MIT License - Astana IT University 2026
```

**Commit message:** `docs(readme): add comprehensive project documentation`

---

#### **File: `.gitignore`**
```
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
cinema-system

# Test binary
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment variables
.env

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db