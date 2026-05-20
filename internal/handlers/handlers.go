package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sabrodigan/bookings-app/internal/config"
	"github.com/sabrodigan/bookings-app/internal/forms"
	"github.com/sabrodigan/bookings-app/internal/helpers"
	"github.com/sabrodigan/bookings-app/internal/mail"
	"github.com/sabrodigan/bookings-app/internal/models"
	"github.com/sabrodigan/bookings-app/internal/render"
	"github.com/sabrodigan/bookings-app/internal/repository"
	"log"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db repository.DatabaseRepo) *Repository {
	return &Repository{
		App: a,
		DB:  db,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	log.Println("Home page hit")
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	log.Println("About page hit")

	// send data to the template
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	log.Println("Reservation page hit")
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	log.Println("PostReservation page hit")
	err := r.ParseForm()
	// err = errors.New("this is as debugging error message for central error handling")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	// using the required to validate entry
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		var data = make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	
	resID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	
	m.App.Session.Put(r.Context(), "reservation_id", resID)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	log.Println("Generals page hit")
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	log.Println("Majors page hit")
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	log.Println("Availability page hit")
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	log.Println("PostAvailability page hit")
	helpers.TypeWriter("PostAvailability page hit", 50)
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	RoomID    string `json:"room_id,omitempty"`
	RoomName  string `json:"room_name,omitempty"`
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	out, err := json.MarshalIndent(v, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// AvailabilityJSON handles requests for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")
	roomID := r.Form.Get("room_id")
	roomName := r.Form.Get("room_name")

	if start == "" || end == "" {
		writeJSON(w, jsonResponse{
			OK:      false,
			Message: "Please choose arrival and departure dates.",
		})
		return
	}

	if start >= end {
		writeJSON(w, jsonResponse{
			OK:      false,
			Message: "Departure must be after arrival.",
		})
		return
	}

	writeJSON(w, jsonResponse{
		OK:        true,
		Message:   "Available!",
		StartDate: start,
		EndDate:   end,
		RoomID:    roomID,
		RoomName:  roomName,
	})
}

// PostBookNow creates a booking from email + dates and sends a confirmation email
func (m *Repository) PostBookNow(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	email := r.Form.Get("email")
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	roomID := r.Form.Get("room_id")
	roomName := r.Form.Get("room_name")

	form := forms.New(r.PostForm)
	form.Required("email", "start", "end")
	form.IsEmail("email")

	if start == "" || end == "" || start >= end {
		form.Errors.Add("", "Invalid dates selected")
	}

	if !form.Valid() {
		writeJSON(w, jsonResponse{
			OK:      false,
			Message: "Please provide a valid email and dates.",
		})
		return
	}

	reservation := models.Reservation{
		FirstName: "Guest",
		LastName:  "Booking",
		Email:     email,
		Phone:     "—",
		StartDate: start,
		EndDate:   end,
		RoomID:    roomID,
		RoomName:  roomName,
	}

	resID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if err := mail.SendBookingConfirmation(reservation); err != nil {
		m.App.ErrorLog.Println("booking confirmation email:", err)
		writeJSON(w, jsonResponse{
			OK:      false,
			Message: "Booking saved but we could not send the confirmation email. Please contact us.",
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation_id", resID)

	writeJSON(w, jsonResponse{
		OK:        true,
		Message:   fmt.Sprintf("Confirmation sent to %s", email),
		StartDate: start,
		EndDate:   end,
		RoomID:    roomID,
		RoomName:  roomName,
	})
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	log.Println("Contact page hit")
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservationID, ok := m.App.Session.Get(r.Context(), "reservation_id").(string)
	if !ok {
		//log.Println("cannot find session data for reservation")
		m.App.ErrorLog.Println("cannot find session data for reservation")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	
	reservation, err := m.DB.GetReservationByID(reservationID)
	if err != nil {
		m.App.ErrorLog.Println("cannot fetch reservation from database")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from database")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
	m.App.Session.Remove(r.Context(), "reservation_id")
}
func (m *Repository) Error(w http.ResponseWriter, r *http.Request) {
	log.Println("Error page has been called")
	// render.RenderTemplate(w, r, "error.page.tmpl", &models.TemplateData{})
}
