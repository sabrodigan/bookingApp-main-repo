package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	data := url.Values{}
	form := New(data)

	if form == nil {
		t.Error("New() should return a non-nil *Form")
	}
}

func TestForm_Valid(t *testing.T) {
	// Create a request with no form data
	r := httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = url.Values{}

	form := New(r.PostForm)

	// A new form with no errors should be valid
	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}

	// Add an error to the form
	form.Errors.Add("test", "This is an error")

	// Form should now be invalid
	isValid = form.Valid()
	if isValid {
		t.Error("got valid when should have been invalid")
	}
}

func TestForm_Required(t *testing.T) {
	// Create a form with empty fields
	r := httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = url.Values{}
	form := New(r.PostForm)

	// Test with a required field that's not in the form
	form.Required("a")
	if form.Valid() {
		t.Error("form shows valid when required field is missing")
	}

	// Add some fields to the form
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "")
	postedData.Add("c", "c")

	form = New(postedData)

	// Test with multiple fields, one of which is empty
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required field b is empty")
	}

	// Check the error message
	if form.Errors.Get("b") != "This field cannot be blank" {
		t.Error("did not get correct error message for empty field")
	}
}

func TestForm_Has(t *testing.T) {
	// Create a request with form data
	postedData := url.Values{}
	postedData.Add("a", "a")

	r := httptest.NewRequest("POST", "/whatever", nil)
	r.Form = postedData

	form := New(postedData)

	// Test with a field that exists
	if !form.Has("a") {
		t.Error("shows form does not have field when it does")
	}

	// Test with a field that doesn't exist
	if form.Has("b") {
		t.Error("shows form has field when it doesn't")
	}
}

func TestForm_MinLength(t *testing.T) {
	// Create a request with form data
	postedData := url.Values{}
	postedData.Add("name", "John")

	r := httptest.NewRequest("POST", "/whatever", nil)
	r.Form = postedData

	form := New(postedData)

	// Test with a field that meets the minimum length
	if !form.MinLength("name", 3) {
		t.Error("shows min length of 3 not met when it is")
	}

	// Test with a field that doesn't meet the minimum length
	if form.MinLength("name", 5) {
		t.Error("shows min length of 5 met when it is not")
	}

	// Check the error message
	if form.Errors.Get("name") != "This field must be at least 5 characters long" {
		t.Error("did not get correct error message for min length")
	}
}

func TestForm_IsEmail(t *testing.T) {
	// Create a form with an email field
	postedData := url.Values{}
	postedData.Add("email", "me@here.com")
	form := New(postedData)

	// Test with a valid email
	if !form.IsEmail("email") {
		t.Error("shows invalid email for valid email address")
	}

	// Test with an invalid email
	postedData = url.Values{}
	postedData.Add("email", "invalid-email")
	form = New(postedData)

	if form.IsEmail("email") {
		t.Error("shows valid email for invalid email address")
	}

	// Check the error message
	if form.Errors.Get("email") != "Invalid email address" {
		t.Error("did not get correct error message for invalid email")
	}
}
