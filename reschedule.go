package main

import (
	"fmt"
	"net/http"
	"time"
)

// handleRescheduleBooking processes rescheduling requests from admin
func handleRescheduleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/admin/bookings", http.StatusSeeOther)
		return
	}

	bookingID := r.FormValue("booking_id")
	newDate := r.FormValue("date")
	newSlot := r.FormValue("time_slot")

	// Validate inputs
	if bookingID == "" || newDate == "" || newSlot == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Update the global bookings slice
	for i, b := range bookings {
		if b.ID == bookingID {
			bookings[i].Date = newDate
			bookings[i].TimeSlot = newSlot
			bookings[i].Status = "rescheduled"
			fmt.Printf("Booking %s rescheduled to %s at %s\n", bookingID, newDate, newSlot)
			break
		}
	}

	http.Redirect(w, r, "/admin/bookings", http.StatusSeeOther)
}
