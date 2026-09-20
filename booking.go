package main

import (
	"fmt"
	"net/http"
	"html/template"
	"strings"
	"time"
)

// booking types
type Booking struct {
	ID             string
	CustomerName   string
	Phone          string
	Email          string
	Service        string
	Date           string
	TimeSlot       string
	Notes          string
	Status         string
	CreatedAt      string
}

// store bookings in memory
var bookings []Booking

// availbale services
var services = map[string]ServiceInfo{
	"consultation": {
		Name:     "consultation",
		Duration: "1 hour 30 minutes",
		Price:    "#10,000",
		Slots: []string{
			"09:00 AM",
			"10:30 AM",
			"12:00 PM",
			"01:30 PM",
			"03:00 PM",
			"04:30 PM",
		},
	},
	"fitting": {
		Name:  "Dress Fitting",
		Duration: "2 hours",
		Price:     "#15,000",
		Slots:  []string{
			"09:00 AM",
			"11:00 AM",
			"01:00 AM",
			"03:00 PM",
		},   
	},
}

type ServiceInfo struct {
	Name      string
	Duration  string
	Price     string
	Slots   []string
}

//booking page data
type BookingData struct {
	Services  map[string]ServiceInfo
	Bookings  []Booking
	Booking   Booking
	Slots     []string
	Error     string
	Success   string
	UserName  string 
}

// render booking template
func renderBoookingTemplate(w http.ResponseWriter, tmp string, data BookingData) {
	t, err := template.ParseFiles(
		"templates/Booking/base.html",
		"templates/booking/"+tmp+".html",
	)
	if err != nil {
		fmt.Println("Booking template error:", err)
		http.Error(w, err.Error(), 500)
		return
	}
    t.ExecuteTemplate(w, "booking-base", data)
}


// render booking template
func renderBookingTemplate(w http.ResponseWriter, tmpl string, data BookingData) {
    t, err := template.ParseFiles(
        "templates/booking/base.html",
        "templates/booking/"+tmpl+".html",
    )
    if err != nil {
        fmt.Println("Booking template error:", err)
        http.Error(w, err.Error(), 500)
        return
    }
    t.ExecuteTemplate(w, "booking-base", data)
}

// check if day is valid working day
func isWorkingDay(date string) bool {
    t, err := time.Parse("2006-01-02", date)
    if err != nil {
        return false
    }

    // Sunday is not a working day
    if t.Weekday() == time.Sunday {
        return false
    }

    return true
}

// check if date is in the past
func isPastDate(date string) bool {
    t, err := time.Parse("2006-01-02", date)
    if err != nil {
        return true
    }
    return t.Before(time.Now().Truncate(24 * time.Hour))
}

// count bookings for a specific date
func countBookingsForDate(date string) int {
    count := 0
    for _, b := range bookings {
        if b.Date == date && b.Status != "cancelled" {
            count++
        }
    }
    return count
}

// check if slot is already taken
func isSlotTaken(date, timeSlot, service string) bool {
    for _, b := range bookings {
        if b.Date == date &&
            b.TimeSlot == timeSlot &&
            b.Service == service &&
            b.Status != "cancelled" {
            return true
        }
    }
    return false
}

// get available slots for a date and service
func getAvailableSlots(date, service string) []string {
    serviceInfo, exists := services[service]
    if !exists {
        return []string{}
    }

    available := []string{}
    for _, slot := range serviceInfo.Slots {
        if !isSlotTaken(date, slot, service) {
            available = append(available, slot)
        }
    }
    return available
}

// generate booking ID
func generateBookingID() string {
    return fmt.Sprintf("ITUA-%d", len(bookings)+1)
}

// GET /booking
func handleBookingHome(w http.ResponseWriter, r *http.Request) {
    renderBookingTemplate(w, "home", BookingData{
        Services: services,
    })
}

// GET /booking/new
// POST /booking/new
func handleNewBooking(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {

        // read form values
        customerName := r.FormValue("name")
        phone        := r.FormValue("phone")
        email        := r.FormValue("email")
        service      := r.FormValue("service")
        date         := r.FormValue("date")
        timeSlot     := r.FormValue("time_slot")
        notes        := r.FormValue("notes")

        // VALIDATE
        if customerName == "" || phone == "" || service == "" ||
            date == "" || timeSlot == "" {
            renderBookingTemplate(w, "new", BookingData{
                Services: services,
                Error:    "Please fill in all required fields",
            })
            return
        }

        // check working day
        if !isWorkingDay(date) {
            renderBookingTemplate(w, "new", BookingData{
                Services: services,
                Error:    "We do not work on Sundays. Please choose Monday to Saturday",
            })
            return
        }

        // check past date
        if isPastDate(date) {
            renderBookingTemplate(w, "new", BookingData{
                Services: services,
                Error:    "Please choose a future date",
            })
            return
        }

        // check max bookings per day
        if countBookingsForDate(date) >= 10 {
            renderBookingTemplate(w, "new", BookingData{
                Services: services,
                Error:    "Sorry this date is fully booked. Please choose another date",
            })
            return
        }

        // check slot available
        if isSlotTaken(date, timeSlot, service) {
            renderBookingTemplate(w, "new", BookingData{
                Services: services,
                Error:    "Sorry this time slot is taken. Please choose another time",
            })
            return
        }

        // create booking
        booking := Booking{
            ID:           generateBookingID(),
            CustomerName: customerName,
            Phone:        phone,
            Email:        email,
            Service:      service,
            Date:         date,
            TimeSlot:     timeSlot,
            Notes:        notes,
            Status:       "pending",
            CreatedAt:    time.Now().Format("02 Jan 2006 15:04"),
        }

        // save booking
        bookings = append(bookings, booking)

        fmt.Println("New booking:", booking.ID, "from:", customerName)

        // show confirmation
        renderBookingTemplate(w, "confirm", BookingData{
            Booking:  booking,
            Services: services,
            Success:  "Your booking has been received!",
        })
        return
    }

    // show booking form
    renderBookingTemplate(w, "new", BookingData{
        Services: services,
    })
}

// GET /booking/slots?date=2024-01-01&service=consultation
func handleGetSlots(w http.ResponseWriter, r *http.Request) {
    date    := r.URL.Query().Get("date")
    service := r.URL.Query().Get("service")

    if date == "" || service == "" {
        http.Error(w, "missing date or service", 400)
        return
    }

    slots := getAvailableSlots(date, service)

    // return slots as HTML options
    var html strings.Builder
    for _, slot := range slots {
        html.WriteString(fmt.Sprintf(
            `<option value="%s">%s</option>`, slot, slot,
        ))
    }

    if len(slots) == 0 {
        html.WriteString(`<option value="">No slots available</option>`)
    }

    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, html.String())
}

// GET /booking/reschedule?id=ITUA-1
// POST /booking/reschedule
func handleReschedule(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        bookingID := r.FormValue("booking_id")
        newDate   := r.FormValue("date")
        newSlot   := r.FormValue("time_slot")

        // validate
        if newDate == "" || newSlot == "" {
            renderBookingTemplate(w, "reschedule", BookingData{
                Error: "Please select a new date and time",
            })
            return
        }

        // check working day
        if !isWorkingDay(newDate) {
            renderBookingTemplate(w, "reschedule", BookingData{
                Error: "We do not work on Sundays",
            })
            return
        }

        // check past date
        if isPastDate(newDate) {
            renderBookingTemplate(w, "reschedule", BookingData{
                Error: "Please choose a future date",
            })
            return
        }

        // find and update booking
        for i, b := range bookings {
            if b.ID == bookingID {
                // check new slot available
                if isSlotTaken(newDate, newSlot, b.Service) {
                    renderBookingTemplate(w, "reschedule", BookingData{
                        Booking: b,
                        Error:   "That slot is taken. Please choose another",
                    })
                    return
                }

                // update booking
                bookings[i].Date     = newDate
                bookings[i].TimeSlot = newSlot
                bookings[i].Status   = "rescheduled"

                fmt.Println("Booking rescheduled:", bookingID)

                renderBookingTemplate(w, "confirm", BookingData{
                    Booking: bookings[i],
                    Success: "Your booking has been rescheduled!",
                })
                return
            }
        }

        http.Error(w, "Booking not found", 404)
        return
    }

    // show reschedule form
    bookingID := r.URL.Query().Get("id")
    for _, b := range bookings {
        if b.ID == bookingID {
            renderBookingTemplate(w, "reschedule", BookingData{
                Booking:  b,
                Services: services,
            })
            return
        }
    }

    http.Error(w, "Booking not found", 404)
}

// admin bookings handler
func handleAdminBookings(w http.ResponseWriter, r *http.Request) {
    if !requireAuth(w, r) {
        return
    }

    session, _ := getSession(r)

    renderAdminTemplate(w, "bookings", AdminData{
        Username: session.Username,
        Bookings: bookings,
        Stats: AdminStats{
            TotalBookings: len(bookings),
        },
    })
}

// POST /admin/bookings/approve
func handleApproveBooking(w http.ResponseWriter, r *http.Request) {
    if !requireAuth(w, r) {
        return
    }

    bookingID := r.FormValue("booking_id")
    for i, b := range bookings {
        if b.ID == bookingID {
            bookings[i].Status = "approved"
            fmt.Println("Booking approved:", bookingID)
            break
        }
    }

    http.Redirect(w, r, "/admin/bookings", http.StatusSeeOther)
}

// POST /admin/bookings/cancel
func handleCancelBooking(w http.ResponseWriter, r *http.Request) {
    if !requireAuth(w, r) {
        return
    }

    bookingID := r.FormValue("booking_id")
    for i, b := range bookings {
        if b.ID == bookingID {
            bookings[i].Status = "cancelled"
            fmt.Println("Booking cancelled:", bookingID)
            break
        }
    }

    http.Redirect(w, r, "/admin/bookings", http.StatusSeeOther)
}