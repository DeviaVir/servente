package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DeviaVir/servente/pkg/models"
	"github.com/DeviaVir/servente/pkg/models/mysql"

	"github.com/golangcollege/sessions"
	"github.com/namsral/flag"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	debug    bool
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	services interface {
		Insert(*models.Organization, string, string, string, []*models.ServiceAttribute, int, string) (int, error)
		Get(*models.Organization, int) (*models.Service, error)
		Latest(*models.Organization, int, int) ([]*models.Service, error)
		Update(*models.Service) (int, error)
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
		UpdateAttribute(*models.Setting, string) (*models.OrganizationAttribute, error)
		UpdateSetting(*models.Setting) (*models.Setting, error)
		Get(string) (*models.Organization, error)
		GetSettings(*models.Organization) ([]*models.Setting, error)
		GetAttributes(*models.Organization) ([]*models.OrganizationAttribute, error)
	}
}

func main() {
	debug := flag.Bool("debug", false, "Enable debug stack traces shown to users")
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "servente:servente@/servente", "MySQL data source name")
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
	db.AutoMigrate(&models.Organization{}, &models.User{}, &models.Service{}, &models.OrganizationAttribute{}, &models.ServiceAttribute{}, &models.Setting{}, &models.AuditLog{})

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

	var tlsConfig *tls.Config
	if fileExists(*tlsCertPath) && fileExists(*tlsKeyPath) {
		cert, err := tls.LoadX509KeyPair(*tlsCertPath, *tlsKeyPath)
		if err != nil {
			errorLog.Fatal(err)
		}

		tlsConfig = &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
			Certificates:             []tls.Certificate{cert},
		}
	} else {
		infoLog.Printf("No TLS certs found, generating self-signed certs.")
		tlsConfig, err = generateTLSCerts()
		if err != nil {
			errorLog.Fatal(err)
		}
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
	err = srv.ListenAndServeTLS("", "")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*gorm.DB, *sql.DB, error) {
	// check if DSN contains parseTime=true (we expect it)
	if !strings.Contains(dsn, "parseTime=true") {
		if strings.Contains(dsn, "?") {
			dsn = fmt.Sprintf("%s&%s", dsn, "parseTime=true")
		} else {
			dsn = fmt.Sprintf("%s?%s", dsn, "parseTime=true")
		}
	}
	// check if DSN contains charset=utf8mb4 (we expect it)
	if !strings.Contains(dsn, "charset=utf8mb4") {
		if strings.Contains(dsn, "?") {
			dsn = fmt.Sprintf("%s&%s", dsn, "charset=utf8mb4")
		} else {
			dsn = fmt.Sprintf("%s?%s", dsn, "charset=utf8mb4")
		}
	}

	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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

func generateTLSCerts() (*tls.Config, error) {
	netInterfaces, err := net.Interfaces()

	if err != nil {
		return nil, fmt.Errorf("failed to load net interfaces: %v", err)
	}

	var ipAddress *net.IP

	for _, netInterface := range netInterfaces {
		ipAddresses, err := netInterface.Addrs()

		if err != nil {
			return nil, fmt.Errorf("failed to load addresses for net interface: %v", err)
		}

		for _, addr := range ipAddresses {
			switch v := addr.(type) {
			case *net.IPNet:
				ipAddress = &v.IP
			case *net.IPAddr:
				ipAddress = &v.IP
			}

			if ipAddress != nil {
				break
			}
		}

		if ipAddress != nil {
			break
		}
	}

	if ipAddress == nil {
		return nil, fmt.Errorf("failed to discover IP address")
	}

	issuer := pkix.Name{CommonName: ipAddress.String()}

	caCertificate := &x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		Subject:               issuer,
		Issuer:                issuer,
		SignatureAlgorithm:    x509.SHA512WithRSA,
		PublicKeyAlgorithm:    x509.ECDSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, 10),
		SubjectKeyId:          []byte{},
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	caCertificateBinary, err := x509.CreateCertificate(rand.Reader, caCertificate, caCertificate, &privateKey.PublicKey, privateKey)

	if err != nil {
		return nil, fmt.Errorf("create cert failed: %v", err)
	}

	caCertificateParsed, _ := x509.ParseCertificate(caCertificateBinary)

	certPool := x509.NewCertPool()
	certPool.AddCert(caCertificateParsed)

	return &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		ServerName:               ipAddress.String(),
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{caCertificateBinary},
			PrivateKey:  privateKey,
		}},
		RootCAs: certPool,
	}, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
