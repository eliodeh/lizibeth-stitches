# Itua Stitches

A Go-based web application built for managing the operations of "Itua Stitches," a fashion brand. The application provides a public-facing website for viewing collections, booking services, and contacting the business, as well as an administrative dashboard for managing content.

## Features
- **Public Website:**
  - View collections, services, and booking availability.
  - Contact form for inquiries.
- **Booking System:**
  - Online service booking and slot management.
  - Rescheduling capabilities for administrators.
- **Admin Dashboard:**
  - Secure authentication (bcrypt password hashing).
  - Management of collections and inventory items.
  - Image upload capabilities for collections.
  - Message board for incoming inquiries.
  - Booking management (approve/cancel/reschedule).

## Installation
1. Ensure you have [Go](https://go.dev/doc/install) installed (Go 1.26+).
2. Clone this repository.
3. Run the application:
   ```bash
   go run main.go
   ```
4. Access the application at `http://localhost:8080`.

## Technical Details
- **Database:** SQLite (using `modernc.org/sqlite` for CGO-free portability).
- **Security:** bcrypt password hashing and secure session management with `crypto/rand`.
- **Infrastructure:** Dockerized with a multi-stage build.

## Project Structure
- `main.go`: Entry point, routing, and server config.
- `database.go`: SQLite initialization and schema setup.
- `auth.go`: Authentication logic, middleware, and session management.
- `admin.go`: Admin dashboard handlers.
- `booking.go`: Booking business logic.
- `upload.go`: Image upload handling.
- `reschedule.go`: Booking rescheduling logic.
- `templates/`: HTML templates.
- `static/`: Publicly served static files.
