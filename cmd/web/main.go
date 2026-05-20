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
	"github.com/sabrodigan/bookings-app/internal/driver"
	"github.com/sabrodigan/bookings-app/internal/handlers"
	"github.com/sabrodigan/bookings-app/internal/helpers"
	"github.com/sabrodigan/bookings-app/internal/models"
	"github.com/sabrodigan/bookings-app/internal/render"
	"github.com/sabrodigan/bookings-app/internal/repository/dbrepo"
	"github.com/joho/godotenv"
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
	_ = godotenv.Load()

	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "postgres"
	}

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

	log.Println("Connecting to database...")
	var dbConn *driver.DB

	switch dbDriver {
	case "mongo":
		mongoURI := os.Getenv("MONGO_URI")
		if mongoURI == "" {
			mongoURI = "mongodb://localhost:27017"
		}
		dbConn, err = driver.ConnectMongo(mongoURI)
		if err != nil {
			log.Fatal("Cannot connect to MongoDB: ", err)
		}
		handlers.NewHandlers(handlers.NewRepo(&app, dbrepo.NewMongoRepo(dbConn.Client, &app)))
		log.Println("Connected to MongoDB!")
	default:
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			host := envOr("POSTGRES_HOST", "127.0.0.1")
			port := envOr("POSTGRES_PORT", "5432")
			user := envOr("POSTGRES_USER", "postgres")
			pass := envOr("POSTGRES_PASSWORD", "postgres")
			dbname := envOr("POSTGRES_DB", "devdb")
			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				host, port, user, pass, dbname)
		}
		dbConn, err = driver.ConnectPostgres(dsn)
		if err != nil {
			log.Fatal("Cannot connect to PostgreSQL: ", err)
		}
		handlers.NewHandlers(handlers.NewRepo(&app, dbrepo.NewPostgresRepo(dbConn.SQL, &app)))
		log.Println("Connected to PostgreSQL!")
	}

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

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// clearScreen clears the terminal screen
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	// It's safe to ignore the error for a non-critical cosmetic function
	_ = cmd.Run()
}
