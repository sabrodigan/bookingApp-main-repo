package render

import (
	"encoding/gob"
	"errors"
	"github.com/alexedwards/scs/v2"
	"github.com/sabrodigan/bookings-app/internal/config"
	"github.com/sabrodigan/bookings-app/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	// change this to true when in production
	testApp.InProduction = false
	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false
	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{
	failWrite bool
}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}
func (tw *myWriter) WriteHeader(i int) {

}
func (tw *myWriter) Write(b []byte) (int, error) {
	if tw.failWrite {
		return 0, errors.New("error writing to response")
	}
	length := len(b)
	return length, nil
}
