// /home/sabrodigan/go/src/bookingApp/cmd/web/main.go
package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/sabrodigan/bookings-app/internal/config"
	"github.com/sabrodigan/bookings-app/internal/handlers"
	"github.com/sabrodigan/bookings-app/internal/helpers"
	"github.com/sabrodigan/bookings-app/internal/models"
	"github.com/sabrodigan/bookings-app/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {
	// this runs the Application server
	clearScreen()

	// run the app from here
	err := run()
	if err != nil {
		log.Fatal("cannot start application: ", err)
	}

	// This message now appears after the server is ready to listen
	helpers.TypeWriter("\n\nListening for requests on port ", 50)
	helpers.TypeWriter(portNumber, 50)
	fmt.Println()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
func run() error {
	// set up the application config
	gob.Register(models.Reservation{})

	// Initialize helpers and the audio system FIRST
	helpers.NewHelpers(&app)
	helpers.InitAudio()

	// Now that audio is initialized, we can print the startup message
	helpers.TypeWriter("Starting the application server... ", 50)
	helpers.TypeWriter("\n\nListening for requests on port ", 50)
	helpers.TypeWriter(portNumber, 50)

	// change this to true when in production
	app.InProduction = false

	// set up the logger
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache: ", err)
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	done := make(chan bool)
	if !app.InProduction {
		helpers.TypeWriter("\nStarting the development server... working\n", 45)
		// I also fixed a bug here: the delay must be a time.Duration
		go helpers.Spinner(50*time.Millisecond, done)
	} else {
		helpers.TypeWriter("\nStarting the production server... working\n", 25)
	}

	return nil
}

// clearScreen clears the terminal screen
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	// It's safe to ignore the error for a non-critical cosmetic function
	_ = cmd.Run()
}
