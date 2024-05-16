package main

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/blue-davinci/gosnip/internal/jsonlog"
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
	debug bool
}

type application struct {
	config         config
	logger         *jsonlog.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	var cfg config
	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	currentpath := getEnvPath(logger)
	logger.PrintInfo("Searching env at: ", map[string]string{"ENV:": currentpath})
	if currentpath != "" {
		err := godotenv.Load(currentpath)
		if err != nil {
			logger.PrintError(err, nil)
		}
	} else {
		return
	}
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("GOSNIP_DB_DSN"), "MySQL DB DSN")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	// show debug
	cfg.debug = *debug
	logger.PrintInfo(fmt.Sprintf("Debug mode: %v", cfg.debug), nil)
	//------------------------------------------- Database
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		os.Exit(0)
	}
	defer db.Close()
	//---------------------------------------------------
	// Initialize new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.PrintError(err, nil)
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
		logger:         logger,
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
	app.logger.PrintInfo("Server is running",
		map[string]string{
			"port":        fmt.Sprintf(":%d", app.config.port),
			"Environment": app.config.env,
		})
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		TLSConfig:    tlsConfig,
	}
	err = s.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if !errors.Is(err, http.ErrServerClosed) {
		app.logger.PrintError(err, nil)
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
func getEnvPath(l *jsonlog.Logger) string {
	dir, err := os.Getwd()
	if err != nil {
		l.PrintError(err, nil)
		return ""
	}
	if strings.Contains(dir, "cmd/web") || strings.Contains(dir, "cmd") {
		return ".env"
	}
	return filepath.Join("cmd", "web", ".env")
}
