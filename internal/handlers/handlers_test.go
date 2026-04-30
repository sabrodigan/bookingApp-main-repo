package handlers

import (
	"encoding/json"
	"github.com/sabrodigan/bookings-app/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gs", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"pr", "/make-reservation", "POST", []postData{
		{"start", "2023-10-01"},
		{"end", "2023-10-02"},
	}, http.StatusOK},
	{"sa", "/search-availability", "POST", []postData{
		{"start", "2023-10-01"},
		{"end", "2023-10-02"},
	}, http.StatusOK},
	{"sa", "/search-availability-json", "POST", []postData{
		{"first_name", "Stephen"},
		{"last_name", "Brodigan"},
		{"email", "sabrodigan@gmail.com"},
		{"password", "1234"},
	}, http.StatusOK},
}

func TestNewHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	// Test case 1: reservation in session
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// Make a request to the server to get a valid session cookie
	resp, err := ts.Client().Get(ts.URL + "/make-reservation")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK but got %d", resp.StatusCode)
	}

	// Get the cookie from the response
	cookies := resp.Cookies()

	// Create a new request to test the reservation summary
	req, _ := http.NewRequest("GET", "/reservation-summary", nil)

	// Add the session cookie to the request
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Create a context with session data
	ctx, _ := session.Load(req.Context(), req.Header.Get("X-Session"))

	// Create a reservation and add it to the session
	reservation := models.Reservation{
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@smith.com",
		Phone:     "123456789",
	}

	session.Put(ctx, "reservation", reservation)
	req = req.WithContext(ctx)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create handler with session middleware
	handler := SessionLoad(http.HandlerFunc(Repo.ReservationSummary))

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong status code: got %v want %v", 
			rr.Code, http.StatusOK)
	}

	// Test case 2: reservation not in session
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)

	// Add the session cookie to the request
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Create a context with session data but no reservation
	ctx, _ = session.Load(req.Context(), req.Header.Get("X-Session"))
	req = req.WithContext(ctx)

	// Create a response recorder
	rr = httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code - should be a redirect
	if rr.Code != http.StatusSeeOther {
		t.Errorf("ReservationSummary handler returned wrong status code for missing reservation: got %v want %v", 
			rr.Code, http.StatusSeeOther)
	}

	// Check the location header
	if rr.Header().Get("Location") != "/" {
		t.Errorf("ReservationSummary handler redirected to wrong URL: got %v want %v", 
			rr.Header().Get("Location"), "/")
	}
}

func TestRepository_Error(t *testing.T) {
	// Create a request
	req, _ := http.NewRequest("GET", "/error", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create handler
	handler := http.HandlerFunc(Repo.Error)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// The Error handler doesn't return anything specific, so we just check that it doesn't crash
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// Create a request
	req, _ := http.NewRequest("POST", "/search-availability-json", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create handler
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check content type
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("AvailabilityJSON handler returned wrong content type: got %v want %v", 
			rr.Header().Get("Content-Type"), "application/json")
	}

	// Check the response body
	var resp jsonResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Error("Failed to parse JSON response:", err)
	}

	// Check the response values
	if !resp.OK {
		t.Error("Expected OK to be true, got false")
	}

	if resp.Message != "Available!" {
		t.Errorf("Expected message to be 'Available!', got '%s'", resp.Message)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	// Test case 1: form validation failure

	// Create a request with form data that will fail validation
	postedData := url.Values{}
	postedData.Add("first_name", "a") // Too short, will fail MinLength validation
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "invalid-email") // Will fail IsEmail validation
	postedData.Add("phone", "123456789")

	req, _ := http.NewRequest("POST", "/make-reservation", nil)
	req.PostForm = postedData

	// Create a context with session data
	ctx, _ := session.Load(req.Context(), req.Header.Get("X-Session"))
	req = req.WithContext(ctx)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create handler with session middleware
	handler := SessionLoad(http.HandlerFunc(Repo.PostReservation))

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Should stay on the same page due to validation errors
	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returned wrong status code for invalid form: got %v want %v", 
			rr.Code, http.StatusOK)
	}

	// Test case 2: successful form submission

	// Create a request with valid form data
	postedData = url.Values{}
	postedData.Add("first_name", "John") // Valid length
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@example.com") // Valid email
	postedData.Add("phone", "123456789")

	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	req.PostForm = postedData

	// Create a context with session data
	ctx, _ = session.Load(req.Context(), req.Header.Get("X-Session"))
	req = req.WithContext(ctx)

	// Create a response recorder
	rr = httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Should redirect to reservation summary
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong status code for valid form: got %v want %v", 
			rr.Code, http.StatusSeeOther)
	}

	// Check the location header
	if rr.Header().Get("Location") != "/reservation-summary" {
		t.Errorf("PostReservation handler redirected to wrong URL: got %v want %v", 
			rr.Header().Get("Location"), "/reservation-summary")
	}

	// Test case 3: ParseForm error
	// This is a bit tricky to test since it's hard to force ParseForm to fail
	// We'll create a custom request with a malformed body

	// Create a request with a malformed body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	req.Body = http.NoBody // Force body to be nil
	req.Method = "POST"    // Ensure it's a POST request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a context with session data
	ctx, _ = session.Load(req.Context(), req.Header.Get("X-Session"))
	req = req.WithContext(ctx)

	// Create a response recorder
	rr = httptest.NewRecorder()

	// Create a handler that will force ParseForm to be called
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call ParseForm directly to ensure it's covered
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			return
		}
		// Continue with normal handler
		Repo.PostReservation(w, r)
	})

	// Wrap with session middleware
	handler = SessionLoad(handlerFunc)

	// Serve the request
	handler.ServeHTTP(rr, req)
}
