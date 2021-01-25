package main

import (
	"crypto/tls"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DeviaVir/servente/pkg/models"
	"github.com/DeviaVir/servente/pkg/models/mysql"

	"github.com/golangcollege/sessions"
	"github.com/namsral/flag"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	debug    bool
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	services interface {
		Insert(string, string, string, []*models.Attribute, int) (int, error)
		Get(int) (*models.Service, error)
		Latest(int) ([]*models.Service, error)
	}
	templateCache map[string]*template.Template
	users         interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
		ChangePassword(int, string, string) error
		Organizations(*models.User) ([]*models.Organization, error)
	}
	organizations interface {
		Insert(*models.User, string, string) (*models.Organization, error)
		Update(*models.User, *models.Organization, string) (*models.Organization, error)
		Get(string) (*models.Organization, error)
		GetSettings(*models.Organization) ([]*models.Setting, error)
	}
}

func main() {
	debug := flag.Bool("debug", false, "Enable debug stack traces shown to users")
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "servente:servente@/servente?charset=utf8mb4&parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6ndh+pPbnzHbS*+9Pk8qGWhtzbpa!ge", "Cookie secret key")
	sessionLifetimeHours := flag.Int("session-lifetime-hours", 12, "Session cookie lifetime")
	tlsCertPath := flag.String("tls-cert-path", "./tls/cert.pem", "TLS certificate path")
	tlsKeyPath := flag.String("tls-key-path", "./tls/key.pem", "TLS key path")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, sqlDB, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer sqlDB.Close()
	db.AutoMigrate(&models.Organization{}, &models.User{}, &models.Service{}, &models.Attribute{}, &models.Setting{}, &models.AuditLog{})

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = time.Duration(*sessionLifetimeHours) * time.Hour

	app := &application{
		debug:         *debug,
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		services:      &mysql.ServiceModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
		organizations: &mysql.OrganizationModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS(*tlsCertPath, *tlsKeyPath)
	errorLog.Fatal(err)
}

func openDB(dsn string) (*gorm.DB, *sql.DB, error) {
	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, nil, err
	}

	return db, sqlDB, nil
}
