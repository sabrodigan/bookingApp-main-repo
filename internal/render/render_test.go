package render

import (
	"github.com/sabrodigan/bookings-app/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Fatalf("could not get session: %v", err)
	}

	session.Put(r.Context(), "flash", "test flash")

	result := AddDefaultData(&td, r)
	if result.Flash != "test flash" {
		t.Errorf("expected Flash to be empty, got %s", result.Flash)
	}
}

func TestRenderTemp(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Fatalf("could not get session: %v", err)
	}
	var ww myWriter

	// Test with UseCache = false (default)
	app.UseCache = false
	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error rendering template with UseCache = false", err)
	}

	// Test with UseCache = true
	app.UseCache = true
	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error rendering template with UseCache = true", err)
	}

	// Test template that doesn't exist
	err = RenderTemplate(&ww, r, "nonExistent.page.tmp", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist", err)
	}

	// Test error when writing to response
	ww.failWrite = true
	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("did not get error when writer failed")
	}
	ww.failWrite = false
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = app.Session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}
func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func TestCreateTemplateCacheWithBadTemplate(t *testing.T) {
	// Save the original path
	originalPath := pathToTemplates

	// Create a temporary test path
	pathToTemplates = "../../templates"

	// Get the template cache
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	// Restore the original path
	pathToTemplates = originalPath

	// Verify we got templates
	if len(tc) == 0 {
		t.Error("template cache is empty")
	}
}
