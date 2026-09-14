package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
)

// admin page data
type AdminData struct {
    Username    string
    Collections []Collection
    Messages    []ContactMessage
    Stats       AdminStats
    Error       string
    Success     string
}

type ContactMessage struct {
    Name     string
    Phone    string
    Email    string
    Interest string
    Message  string
    Date     string
}

type AdminStats struct {
    TotalCollections int
    TotalItems       int
    TotalMessages    int
}

// store messages in memory
var contactMessages []ContactMessage

// render admin template
func renderAdminTemplate(w http.ResponseWriter, tmpl string, data AdminData) {
    t, err := template.ParseFiles(
        "templates/admin/base.html",
        "templates/admin/"+tmpl+".html",
    )
    if err != nil {
        fmt.Println("Admin template error:", err)
        http.Error(w, "Template error", 500)
        return
    }
    t.ExecuteTemplate(w, "admin-base", data)
}

// GET /admin/login
// POST /admin/login
func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
    // if already logged in go to dashboard
    if _, valid := getSession(r); valid {
        http.Redirect(w, r, "/admin", http.StatusSeeOther)
        return
    }

    if r.Method == "POST" {
        username := r.FormValue("username")
        password := r.FormValue("password")

        // check credentials
        if username == AdminUsername && password == AdminPassword {
            // create session
            sessionID := createSession(username)

            // set cookie
            http.SetCookie(w, &http.Cookie{
                Name:     "admin_session",
                Value:    sessionID,
                Path:     "/",
                HttpOnly: true,
                MaxAge:   86400, // 24 hours
            })

            http.Redirect(w, r, "/admin", http.StatusSeeOther)
            return
        }

        // wrong credentials
        t, _ := template.ParseFiles(
            "templates/admin/base.html",
            "templates/admin/login.html",
        )
        t.ExecuteTemplate(w, "admin-base", AdminData{
            Error: "Wrong username or password",
        })
        return
    }

    // show login page
    t, err := template.ParseFiles(
        "templates/admin/base.html",
        "templates/admin/login.html",
    )
    if err != nil {
        http.Error(w, "Template error", 500)
        return
    }
    t.ExecuteTemplate(w, "admin-base", AdminData{})
}

// GET /admin
func handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
    if !requireAuth(w, r) {
        return
    }

    session, _ := getSession(r)
    data := getSiteData()

    // count total items
    totalItems := 0
    for _, c := range data.Collections {
        totalItems += len(c.Items)
    }

    renderAdminTemplate(w, "dashboard", AdminData{
        Username:    session.Username,
        Collections: data.Collections,
        Messages:    contactMessages,
        Stats: AdminStats{
            TotalCollections: len(data.Collections),
            TotalItems:       totalItems,
            TotalMessages:    len(contactMessages),
        },
    })
}

// GET /admin/collections
func handleAdminCollections(w http.ResponseWriter, r *http.Request) {
    if !requireAuth(w, r) {
        return
    }

    session, _ := getSession(r)
    data := getSiteData()

    renderAdminTemplate(w, "collections", AdminData{
        Username:    session.Username,
        Collections: data.Collections,
    })
}

// POST /admin/items/add
func handleAdminAddItem(w http.ResponseWriter, r *http.Request) {
    if !requireAuth(w, r) {
        return
    }

    if r.Method != "POST" {
        http.Redirect(w, r, "/admin/collections", http.StatusSeeOther)
        return
    }

    // read form values
    collectionIndex, _ := strconv.Atoi(r.FormValue("collection_index"))
    name  := r.FormValue("name")
    price := r.FormValue("price")
    image := r.FormValue("image")

    // add item to collection
    if collectionIndex >= 0 && collectionIndex < len(siteCollections) {
        siteCollections[collectionIndex].Items = append(
            siteCollections[collectionIndex].Items,
            Item{
                Name:  name,
                Price: price,
                Image: image,
            },
        )
    }

    fmt.Println("New item added:", name, "to collection", collectionIndex)
    http.Redirect(w, r, "/admin/collections", http.StatusSeeOther)
}

// POST /admin/items/delete
func handleAdminDeleteItem(w http.ResponseWriter, r *http.Request) {
    if !requireAuth(w, r) {
        return
    }

    if r.Method != "POST" {
        http.Redirect(w, r, "/admin/collections", http.StatusSeeOther)
        return
    }

    collectionIndex, _ := strconv.Atoi(r.FormValue("collection_index"))
    itemIndex, _       := strconv.Atoi(r.FormValue("item_index"))

    // delete item from collection
    if collectionIndex >= 0 && collectionIndex < len(siteCollections) {
        items := siteCollections[collectionIndex].Items
        if itemIndex >= 0 && itemIndex < len(items) {
            siteCollections[collectionIndex].Items = append(
                items[:itemIndex],
                items[itemIndex+1:]...,
            )
        }
    }

    http.Redirect(w, r, "/admin/collections", http.StatusSeeOther)
}

// GET /admin/messages
func handleAdminMessages(w http.ResponseWriter, r *http.Request) {
    if !requireAuth(w, r) {
        return
    }

    session, _ := getSession(r)

    renderAdminTemplate(w, "messages", AdminData{
        Username: session.Username,
        Messages: contactMessages,
        Stats: AdminStats{
            TotalMessages: len(contactMessages),
        },
    })
}

// GET /admin/logout
func handleAdminLogout(w http.ResponseWriter, r *http.Request) {
    deleteSession(r)

    // clear cookie
    http.SetCookie(w, &http.Cookie{
        Name:   "admin_session",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    })

    http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}