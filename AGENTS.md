# Project Context
This is a Go-based web application for "Itua Stitches".

## Key Commands
- **Run the server:** `go run main.go`
- **Port:** Uses environment variable `PORT`, defaulting to `:8080`.

## Architecture & Dependencies
- **Tech Stack:** `net/http`, `html/template`.
- **Database:** SQLite via `modernc.org/sqlite` (Pure Go, no CGO).
- **Security:** `golang.org/x/crypto/bcrypt` for passwords; `crypto/rand` for session tokens.
- **File Structure:**
  - `database.go`: SQLite schema and init.
  - `auth.go`: Middleware `requireAuth`, bcrypt hashing, secure sessions.
  - `upload.go`: Image handling.
  - `reschedule.go`: Booking management.
  - `templates/`: Always include `templates/base.html` or `templates/admin/base.html`.

## Operational Gotchas
- The application requires `lizibeth.db` to be initialized.
- Admin routes are protected by the `requireAuth` middleware.
- Image uploads are restricted to 10MB via `MultipartForm`.
- Templates reside in specialized subdirectories (booking, admin).
