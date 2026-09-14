package main
import (
    "crypto/sha256"
    "fmt"
    "net/http"
    "sync"
    "time"
)

// admin credentials
const (
    AdminUsername = "Ejembe Itua"
    AdminPassword = "123456789Eli."
)

// session store
var (
    sessions = map[string]Session{}
    mu       sync.Mutex
)

type Session struct {
    Username  string
    ExpiresAt time.Time
}

// hash password for security
func hashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    return fmt.Sprintf("%x", h.Sum(nil))
}

// create a new session
func createSession(username string) string {
    mu.Lock()
    defer mu.Unlock()

    // generate session id
    sessionID := fmt.Sprintf("%x", sha256.New().Sum(
        []byte(username+time.Now().String()),
    ))

    // store session for 24 hours
    sessions[sessionID] = Session{
        Username:  username,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }

    return sessionID
}

// check if session is valid
func getSession(r *http.Request) (Session, bool) {
    cookie, err := r.Cookie("admin_session")
    if err != nil {
        return Session{}, false
    }

    mu.Lock()
    defer mu.Unlock()

    session, exists := sessions[cookie.Value]
    if !exists {
        return Session{}, false
    }

    // check if session expired
    if time.Now().After(session.ExpiresAt) {
        delete(sessions, cookie.Value)
        return Session{}, false
    }

    return session, true
}

// delete session on logout
func deleteSession(r *http.Request) {
    cookie, err := r.Cookie("admin_session")
    if err != nil {
        return
    }

    mu.Lock()
    defer mu.Unlock()
    delete(sessions, cookie.Value)
}

// check if admin is logged in
func requireAuth(w http.ResponseWriter, r *http.Request) bool {
    _, valid := getSession(r)
    if !valid {
        http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
        return false
    }
    return true
}
