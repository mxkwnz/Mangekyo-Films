# Cinema System - GoLang Backend

A comprehensive cinema booking and management system built with Go, Gin framework, and MongoDB.

## Team Members

- Mukhammedali Kalen
- Shyngys Abdullaev
- Alikhan Orynbasarov

Astana IT University - 2026

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Database Models](#database-models)
- [Authentication](#authentication)
- [Business Logic](#business-logic)
- [Testing](#testing)
- [License](#license)

---

## Features

### Core Functionality
- User Authentication & Authorization - JWT-based auth with role management (User/Admin)
- Movie Catalog - Complete movie management with genres, ratings, and trailers
- Session Scheduling - Smart session scheduling with conflict detection
- Ticket Booking - Multi-seat booking with real-time availability
- Payment System - Integrated payment cards and balance management
- Review System - User reviews with automatic rating calculation
- Hall Management - Multiple hall types (Standard, VIP, IMAX, 3D)

### Advanced Features
- Concurrent Booking - Thread-safe booking with goroutines and mutex locks
- Async Processing - Background rating updates using goroutines
- Smart Pricing - Multiple ticket types (Adult, Student, Kid, Pension)
- Seat Selection - Row and seat number validation
- Refund System - Automatic refunds on booking cancellation
- Secure Payments - CVV hashing with bcrypt

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.21+ |
| Framework | Gin Web Framework |
| Database | MongoDB 5.0+ |
| Authentication | JWT (golang-jwt/jwt/v5) |
| Password Hashing | bcrypt |
| Environment | godotenv |
| Validation | Gin validator |

---

## Project Structure

```
cinema-system/
├── main.go                      # Application entry point
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── .env                         # Environment configuration
├── .gitignore                   # Git ignore rules
│
├── internal/
│   ├── config/
│   │   └── database.go          # MongoDB connection setup
│   │
│   ├── models/                  # Data models
│   │   ├── user.go              # User, roles, auth DTOs
│   │   ├── movie.go             # Movie, genre, review
│   │   ├── hall.go              # Cinema halls
│   │   ├── session.go           # Movie sessions
│   │   ├── ticket.go            # Booking tickets
│   │   ├── payment.go           # Payment records
│   │   ├── payment_card.go      # Payment cards
│   │   ├── payment_code.go      # Promo codes
│   │   └── movie_genre.go       # Movie-genre relations
│   │
│   ├── repositories/            # Database layer
│   │   ├── user_repository.go
│   │   ├── movie_repository.go
│   │   ├── hall_repository.go
│   │   ├── session_repository.go
│   │   ├── ticket_repository.go
│   │   ├── review_repository.go
│   │   ├── genre_repository.go
│   │   ├── movie_genre_repository.go
│   │   ├── payment_repository.go
│   │   └── payment_card_repository.go
│   │
│   ├── services/                # Business logic layer
│   │   ├── auth_service.go      # Authentication & user management
│   │   ├── movie_service.go     # Movie operations
│   │   ├── session_service.go   # Session scheduling
│   │   ├── booking_service.go   # Ticket booking logic
│   │   ├── review_service.go    # Reviews & ratings
│   │   ├── genre_service.go     # Genre management
│   │   ├── payment_service.go   # Payment processing
│   │   ├── payment_card_service.go
│   │   ├── movie_genre_service.go
│   │   └── booking_worker.go    # Background worker
│   │
│   ├── handlers/                # HTTP request handlers
│   │   ├── auth_handler.go      # Auth endpoints
│   │   ├── movie_handler.go     # Movie endpoints
│   │   ├── session_handler.go   # Session endpoints
│   │   ├── booking_handler.go   # Booking endpoints
│   │   ├── review_handler.go    # Review endpoints
│   │   ├── hall_handler.go      # Hall endpoints
│   │   ├── genre_handler.go     # Genre endpoints
│   │   ├── payment_handler.go   # Payment endpoints
│   │   ├── payment_card_handler.go
│   │   └── ticket_handler.go
│   │
│   ├── middleware/              # HTTP middleware
│   │   └── auth_middleware.go   # JWT & role validation
│   │
│   └── routes/
│       └── router.go            # Route definitions
│
└── API_DOCUMENTATION.md         # Complete API reference
```

---

## Installation

### Prerequisites

- Go 1.21 or higher
- MongoDB 5.0 or higher
- Git

### Step 1: Install MongoDB

**Ubuntu/Debian**
```bash
sudo apt update
sudo apt install mongodb -y
sudo systemctl start mongodb
sudo systemctl enable mongodb
```

**macOS**
```bash
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
```

**Windows**
Download and install from MongoDB Download Center: https://www.mongodb.com/try/download/community

### Step 2: Clone Repository

```bash
git clone https://github.com/yourusername/cinema-system.git
cd cinema-system
```

### Step 3: Install Dependencies

```bash
go mod tidy
```

This will download all required Go packages:
- github.com/gin-gonic/gin - Web framework
- go.mongodb.org/mongo-driver - MongoDB driver
- github.com/golang-jwt/jwt/v5 - JWT authentication
- golang.org/x/crypto - Password hashing
- github.com/joho/godotenv - Environment variables

---

## Configuration

### Create Environment File

Create a `.env` file in the project root:

```bash
cp .env.example .env
```

### Environment Variables

Edit `.env` with your configuration:

```env
# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=cinema_db

# Server Configuration
PORT=8080
GIN_MODE=release

# Security
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

**MongoDB Atlas (Cloud) Configuration**

For MongoDB Atlas, use this format:

```env
MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/cinema_db?retryWrites=true&w=majority
```

---

## Running the Application

### Start MongoDB

Ensure MongoDB is running:

```bash
# Ubuntu/Debian
sudo systemctl status mongodb

# macOS
brew services list

# If not running, start it:
sudo systemctl start mongodb                    # Ubuntu/Debian
brew services start mongodb-community           # macOS
```

### Run the Application

```bash
go run main.go
```

You should see:

```
Successfully connected to MongoDB!
Cinema System Server starting on port 8080
[GIN-debug] Listening and serving HTTP on :8080
```

### Test the Server

```bash
curl http://localhost:8080/api/movies
```

---

## API Documentation

### Quick Overview

The API has 50+ endpoints organized into three categories:

**Public Endpoints (No Authentication)**
- User registration and login
- Browse movies and sessions
- View genres and reviews
- Check seat availability

**User Endpoints (Authentication Required)**
- Profile management
- Ticket booking and cancellation
- Review creation and management
- Payment card management
- Balance top-up and payments

**Admin Endpoints (Admin Role Required)**
- Movie, hall, and session management
- Genre management
- View all bookings and payments
- Review moderation

### API Base URL

```
http://localhost:8080/api
```

### Example: Register and Book a Ticket

**1. Register User**
```bash
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

Response:
```json
{
  "user": {
    "id": "...",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "role": "USER",
    "balance": 1000.0
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "registration successful"
}
```

**2. Browse Movies**
```bash
curl http://localhost:8080/api/movies
```

**3. Get Sessions**
```bash
curl http://localhost:8080/api/sessions/movie/{movieId}
```

**4. Book Tickets**
```bash
curl -X POST http://localhost:8080/api/bookings \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "...",
    "seats": [
      {
        "row_number": 5,
        "seat_number": 10,
        "type": "ADULT"
      }
    ]
  }'
```

### Complete API Reference

For detailed API documentation with all endpoints, request/response formats, and examples, see API_DOCUMENTATION.md

---

## Database Models

### User
```go
type User struct {
    ID           ObjectID   // MongoDB ID
    FirstName    string     // User's first name
    LastName     string     // User's last name
    Email        string     // Unique email
    PhoneNumber  string     // Contact number
    PasswordHash string     // bcrypt hashed password
    Role         Role       // GUEST, USER, ADMIN
    Balance      float64    // Account balance (default: 1000.0)
    CreatedAt    time.Time  // Registration date
}
```

### Movie
```go
type Movie struct {
    ID           ObjectID     // MongoDB ID
    Name         string       // Movie title
    AgeRating    string       // G, PG, PG-13, R, etc.
    Duration     int          // Duration in minutes
    Description  string       // Movie description
    PosterURL    string       // Poster image URL
    TrailerURL   string       // Trailer video URL
    AgeLimit     int          // Minimum age
    Rating       float64      // Average rating (0-10)
    Genres       []ObjectID   // Genre IDs
    IsComingSoon bool         // Release status
    CreatedAt    time.Time    // Added date
}
```

### Session
```go
type Session struct {
    ID        ObjectID   // MongoDB ID
    MovieID   ObjectID   // Related movie
    HallID    ObjectID   // Cinema hall
    StartTime time.Time  // Session start
    EndTime   time.Time  // Auto-calculated end
    Price     float64    // Base ticket price
}
```

### Ticket
```go
type Ticket struct {
    ID         ObjectID     // MongoDB ID
    UserID     ObjectID     // Ticket owner
    SessionID  ObjectID     // Related session
    PaymentID  ObjectID     // Payment record
    RowNumber  int          // Seat row
    SeatNumber int          // Seat number
    Type       TicketType   // ADULT, STUDENT, KID, PENSION
    Price      float64      // Final price (after discount)
    Status     TicketStatus // BOOKED, PAID, CANCELLED
    CreatedAt  time.Time    // Booking time
}
```

### Hall
```go
type Hall struct {
    ID          ObjectID // MongoDB ID
    Name        string   // Hall name
    Type        HallType // STANDARD, VIP, IMAX, 3D
    Location    string   // Physical location
    TotalRows   int      // Number of rows
    SeatsPerRow int      // Seats in each row
}
```

---

## Authentication

### JWT Token Authentication (Production)

The system uses JWT tokens for secure authentication:

1. Register/Login to receive a JWT token
2. Include token in Authorization header for protected endpoints:

```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

3. Token expires after 72 hours

### Header-Based Authentication (Testing Only)

For development/testing without JWT:

```bash
X-User-ID: 507f1f77bcf86cd799439011
X-User-Role: USER
```

Warning: Only use header-based auth in development environments.

---

## Business Logic

### Ticket Pricing

The system supports multiple ticket types with automatic price calculation:

| Ticket Type | Price Multiplier | Notes |
|------------|------------------|-------|
| ADULT | 100% | Standard pricing |
| STUDENT | 80% | Requires student verification |
| PENSION | 70% | Senior citizen discount |
| KID | 50% | Not available for 18+ movies |

Example:
- Base price: 2500 tenge
- Adult ticket: 2500 tenge (100%)
- Student ticket: 2000 tenge (80%)
- Kid ticket: 1250 tenge (50%)

### Session Scheduling Rules

1. Future Only - Cannot schedule sessions in the past
2. No Conflicts - Hall must be available (no overlapping sessions)
3. Auto End Time - Calculated based on movie duration
4. Cleanup Time - 15-minute buffer between sessions (optional)

### Booking Validation

Before confirming a booking, the system validates:

- Session exists and is in the future
- Seats are available (not already booked)
- Seats are within hall capacity
- User has sufficient balance
- Age restrictions (no KID tickets for 18+ movies)

### Review System

- Rating Scale: 0-10
- One review per user per movie
- Automatic Rating Update: Movie rating recalculated asynchronously
- Average Calculation: Uses all reviews for the movie

---

## Testing

### Manual Testing

**1. Start the Server**
```bash
go run main.go
```

**2. Test Registration**
```bash
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

**3. Test Login**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

**4. Test Movie Listing**
```bash
curl http://localhost:8080/api/movies
```

### Testing with Postman

1. Import the API endpoints into Postman
2. Create an environment with variables:
   - `base_url`: `http://localhost:8080/api`
   - `token`: (set after login)
3. Use `{{base_url}}` and `{{token}}` in requests

---

## Concurrency Features

The system leverages Go's powerful concurrency primitives:

### 1. Goroutines for Async Operations

**Rating Calculation:**
```go
// Reviews trigger async rating updates
go s.updateMovieRating(context.Background(), review.MovieID)
```

**Batch Processing:**
```go
// Multiple ratings updated concurrently
for _, movieID := range movieIDs {
    wg.Add(1)
    go func(id ObjectID) {
        defer wg.Done()
        updateRating(id)
    }(movieID)
}
```

### 2. Mutex Locks for Critical Sections

**Booking Protection:**
```go
s.mu.Lock()
defer s.mu.Unlock()

// Check seat availability
// Create booking
// Deduct balance
```

This prevents race conditions when multiple users book the same seat simultaneously.

### 3. Worker Queues

**Background Processing:**
```go
type BookingWorker struct {
    Queue chan string
}

// Notifications, emails, etc.
worker.Queue <- "Ticket booked successfully"
```

---

## API Endpoints Summary

### Public Endpoints

**Authentication**
- POST /api/auth/register - Register new user
- POST /api/auth/login - Login user

**Movies**
- GET /api/movies - Get all movies
- GET /api/movies/:id - Get movie details

**Sessions**
- GET /api/sessions/upcoming - Get upcoming sessions
- GET /api/sessions/movie/:movieId - Get sessions for movie
- GET /api/sessions/:id - Get session details
- GET /api/sessions/:id/booked-seats - Get booked seats

**Browse**
- GET /api/halls/:id - Get hall details
- GET /api/genres - Get all genres
- GET /api/reviews/movie/:movieId - Get movie reviews

### User Endpoints (Requires Authentication)

**Profile**
- GET /api/auth/me - Get current user profile
- PUT /api/auth/me - Update profile

**Bookings**
- POST /api/bookings - Book tickets
- GET /api/bookings/my - Get my tickets
- DELETE /api/bookings/:id - Cancel booking

**Reviews**
- POST /api/reviews - Create review
- GET /api/reviews/my - Get my reviews
- PUT /api/reviews/:id - Update review
- DELETE /api/reviews/:id - Delete review

**Payment Cards**
- POST /api/payment-cards - Add payment card
- GET /api/payment-cards - Get my cards
- GET /api/payment-cards/:id - Get card details
- DELETE /api/payment-cards/:id - Delete card

**Payments**
- POST /api/payments/topup - Top up balance
- GET /api/payments - Get my payments
- GET /api/payments/:id - Get payment details
- POST /api/payments/:id/refund - Refund payment

### Admin Endpoints (Requires Admin Role)

**Movies**
- POST /api/admin/movies - Create movie
- PUT /api/admin/movies/:id - Update movie
- DELETE /api/admin/movies/:id - Delete movie

**Halls**
- POST /api/admin/halls - Create hall
- GET /api/admin/halls - Get all halls
- PUT /api/admin/halls/:id - Update hall
- DELETE /api/admin/halls/:id - Delete hall

**Sessions**
- POST /api/admin/sessions - Create session
- PUT /api/admin/sessions/:id - Update session
- DELETE /api/admin/sessions/:id - Delete session

**Bookings**
- GET /api/admin/bookings - View all bookings
- GET /api/admin/bookings/session/:sessionId - Get session bookings

**Genres**
- POST /api/admin/genres - Create genre
- PUT /api/admin/genres/:id - Update genre
- DELETE /api/admin/genres/:id - Delete genre

**Reviews**
- DELETE /api/admin/reviews/:id - Delete review

**Payments**
- GET /api/admin/payments - View all payments
- GET /api/admin/payments/user/:userId - Get user payments
- GET /api/admin/payment-cards/user/:userId - Get user cards

---

## License

This project is licensed under the MIT License.

```
MIT License

Copyright (c) 2026 Astana IT University

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## Contact

**Team Members**
- Mukhammedali Kalen - Project Lead
- Shyngys Abdullaev - Backend Developer
- Alikhan Orynbasarov - Backend Developer

**Institution**
Astana IT University
Nur-Sultan, Kazakhstan
2026

**Support**
- Issues: GitHub Issues
- Documentation: API_DOCUMENTATION.md

---

Built with passion by Astana IT University Students
