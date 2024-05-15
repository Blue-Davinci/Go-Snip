package main

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/blue-davinci/gosnip/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type logger struct {
	log slog.Logger
}

type application struct {
	config         config
	logger         *logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	var cfg config

	options := slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewJSONHandler(os.Stdout, &options)
	mylogger := slog.New(handler)
	logger := logger{log: *mylogger}
	currentpath := getEnvPath(logger.log)
	logger.log.Info("Searching env at: ", slog.String("ENV:", currentpath))
	if currentpath != "" {
		err := godotenv.Load(currentpath)
		if err != nil {
			logger.log.Error(err.Error())
		}
	} else {
		return
	}
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("GOSNIP_DB_DSN"), "MySQL DB DSN")
	flag.Parse()
	//------------------------------------------- Database
	db, err := openDB(cfg)
	if err != nil {
		logger.log.Error(err.Error())
		os.Exit(0)
	}
	defer db.Close()
	//---------------------------------------------------
	// Initialize new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.log.Error(err.Error())
		os.Exit(0)
	}
	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()
	// Use the scs.New() function to initialize a new session manager. Then we
	// configure it to use our MySQL database as the session store, and set an
	// expiry after 12 hours
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Make sure that the Secure attribute is set on our session cookies.
	sessionManager.Cookie.Secure = true

	// Initialize our app with all the dependancies
	app := &application{
		config:         cfg,
		logger:         &logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	// Initialize a tls.Config struct to hold the non-default TLS settings for the server
	// Changing the preferences for elliptic curves with assembly implementations to be used.

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	app.logger.log.Info(fmt.Sprintf("Server is running on port: %d || Environment: %s", app.config.port, app.config.env))
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      app.routes(),
		ErrorLog:     log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		TLSConfig:    tlsConfig,
	}
	err = s.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if !errors.Is(err, http.ErrServerClosed) {
		app.logger.log.Error(err.Error())
		return
	}
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// getEnvPath() method returns the path to the .env file based on the current working directory.
func getEnvPath(l slog.Logger) string {
	dir, err := os.Getwd()
	if err != nil {
		l.Error(err.Error())
		return ""
	}
	if strings.Contains(dir, "cmd/web") || strings.Contains(dir, "cmd") {
		return ".env"
	}
	return filepath.Join("cmd", "web", ".env")
}
