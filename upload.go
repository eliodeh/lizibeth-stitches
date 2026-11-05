package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// handleImageUpload processes file uploads from the admin dashboard
func handleImageUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/admin/collections", http.StatusSeeOther)
		return
	}

	// 10MB limit
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error retrieving file:", err)
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure destination directory exists
	uploadPath := "static/images/"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, 0755)
	}

	// Create new file
	dst, err := os.Create(filepath.Join(uploadPath, handler.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded content
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Uploaded file: %s\n", handler.Filename)
	http.Redirect(w, r, "/admin/collections?success=1", http.StatusSeeOther)
}
