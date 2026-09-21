package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

//global collections store
var siteCollections []Collection

func init() {
	siteCollections = getSiteData().Collections
}

// all page data
type PageData struct {
	Title       string
	Collections []Collection
	CEO         CEOProfile
	Contact     ContactInfo
}

type Collection struct {
	Name        string
	Description string
	Image       string
	Items       []Item
}

type Item struct {
	Name  string
	Price string
	Image string
}

type CEOProfile struct {
	Name    string
	Title   string
	Message string
	Image   string
}

type ContactInfo struct {
	Address string
	Phone1  string
	Phone2  string
	Email   string
	MapLink string
}

// website data

func getSiteData() PageData {
	return PageData{
		Title: "Itua Stitches",
		Collections: []Collection{
			{
				Name:        "Bridal Collection",
				Description: "Exquisite bridal wear for your perfect day",
				Image:       "/static/images/collections/bridal.jpg",
				Items: []Item{
					{Name: "studded lace", Price: "₦600,000", Image: "/static/images/collections/bridal1.jpg"},
					{Name: "receiption gown", Price: "₦450,000", Image: "/static/images/collections/bridal2.jpg"},
					{Name: "ankara dresses", Price: "₦330,000", Image: "/static/images/collections/bridal3.jpg"},
				},
			},
			{
				Name:        "Casual Collection",
				Description: "Elegant everyday wear for the modern woman",
				Image:       "/static/images/collections/casual.jpg",
				Items: []Item{
					{Name: "office Chic", Price: "₦45,000", Image: "/static/images/collections/casual1.jpg"},
					{Name: "Urban beauty", Price: "₦55,000", Image: "/static/images/collections/casual2.jpg"},
					{Name: "City bloom", Price: "₦50,000", Image: "/static/images/collections/casual3.jpg"},
				},
			},
			{
				Name:        "Ready to Wears",
				Description: "Stunning ready to wear gowns for every occasion",
				Image:       "/static/images/collections/evening.jpg",
				Items: []Item{
					{Name: "Midnight Gold", Price: "₦95,000", Image: "/static/images/collections/evening1.jpg"},
					{Name: "Scarlet Night", Price: "₦85,000", Image: "/static/images/collections/evening2.jpg"},
					{Name: "Diamond Dusk", Price: "₦110,000", Image: "/static/images/collections/evening3.jpg"},
				},
			},
			{
				Name:        "Accessories",
				Description: "Complete your look with our luxury accessories",
				Image:       "/static/images/collections/accessories.jpg",
				Items: []Item{
					{Name: "braclets", Price: "₦25,000", Image: "/static/images/collections/acc1.jpg"},
					{Name: "Necklace", Price: "₦35,000", Image: "/static/images/collections/acc2.jpg"},
					{Name: "faciliator", Price: "₦15,000", Image: "/static/images/collections/acc3.jpg"},
				},
			},
		},
		CEO: CEOProfile{
			Name:    "Ejembe Itua",
			Title:   "Founder & CEO",
			Message: "At Itua Stitches, we believe every woman deserves to feel elegant. Our designs are crafted with passion, precision and love for being woman.",
			Image:   "/static/images/ceo.jpg",
		},
		Contact: ContactInfo{
			Address: "Gwarinpa. Abuja",
			Phone1:  "+2349093348066",
			Email:   "ituastitches@gmail.com",
			MapLink: "https://maps.google.com/?q=Gwarinpa+Abuja",
		},
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	t, err := template.ParseFiles(
		"templates/base.html",
		"templates/"+tmpl+".html",
	)
	if err != nil {

		fmt.Println("Template error:", err)
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println("Executing template:", tmpl)

	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Println("Execute error:", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", getSiteData())
}

func handleCollections(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "collections", getSiteData())
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about", getSiteData())
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact", getSiteData())
}

func handleContactForm(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Redirect(w, r, "/contact", http.StatusSeeOther)
        return
    }

	// save message
	contactMessages = append(contactMessages, ContactMessage{
        Name:     r.FormValue("name"),
        Phone:    r.FormValue("phone"),
        Email:    r.FormValue("email"),
        Interest: r.FormValue("interest"),
        Message:  r.FormValue("message"),
        Date:     time.Now().Format("02 Jan 2006 15:04"),
    })

    // print to terminal
    fmt.Println("New message from:", r.FormValue("name"))

    // redirect ONCE only
    http.Redirect(w, r, "/contact", http.StatusSeeOther)
}

func main() {
	// serve static files
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	//  puplic routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/collections", handleCollections)
	http.HandleFunc("/about", handleAbout)
	http.HandleFunc("/contact", handleContact)
	http.HandleFunc("/contact/send", handleContactForm)
    
	// booking routes
    http.HandleFunc("/booking", handleBookingHome)
	http.HandleFunc("/booking/new", handleNewBooking)
	http.HandleFunc("/booking/slots", handleGetSlots)
	http.HandleFunc("/bookingreschedule", handleReschedule)

    // admin routes
	http.HandleFunc("/admin/login", handleAdminLogin)
	http.HandleFunc("/admin/", handleAdminDashboard)
	http.HandleFunc("/admin/collections", handleAdminCollections)
	http.HandleFunc("/admin/items/add", handleAdminAddItem)
	http.HandleFunc("/admin/items/delete", handleAdminDeleteItem)
	http.HandleFunc("/admin/messages", handleAdminMessages)
	http.HandleFunc("/admin/bookings", handleAdminBookings)
	http.HandleFunc("/admin/bookings/approve", handleApproveBooking)
	http.HandleFunc("/admin/bookingscancel", handleCancelBooking)
	http.HandleFunc("/admin/logout", handleAdminLogout)

	fmt.Println("itua Stitches running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

	port := os.Getenv("PORT")
	if port == "" {
	    port = "8080"
	}
	fmt.Println("Itua Stitches running on port:", port)
	http.ListenAndServe(":"+port, nil)

}

