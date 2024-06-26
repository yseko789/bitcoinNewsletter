package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/yseko789/bitcoinNewsletter/internal/data"
	"github.com/yseko789/bitcoinNewsletter/internal/mailer"
	pb "github.com/yseko789/bitcoinNewsletter/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
		// maxOpenConns int
		// maxIdleConns int
		// maxIdleTime  time.Duration
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

var grpcClient pb.SummaryServiceClient

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// connect to grpc server
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("error")
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	var addr string = os.Getenv("grpcAddress")
	conn, err := grpc.Dial(addr, grpc.WithAuthority(addr), grpc.WithTransportCredentials(cred))

	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v\n", err)
	}
	defer conn.Close()
	grpcClient = pb.NewSummaryServiceClient(conn)

	// startup server
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DSN"), "PostgreSQL DSN")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.gmail.com", "SMTP hsot")
	flag.IntVar(&cfg.smtp.port, "smpt-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("gmail"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("gmailPassword"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "yseko789 <no-reply@github.com/yseko789>", "SMTP sender")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	cleanup, err := pgxv4.RegisterDriver("cloudsql-postgres", cloudsqlconn.WithIAMAuthN())
	if err != nil {
		return nil, err
	}
	defer cleanup()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("connectionName"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))
	db, err := sql.Open(
		"cloudsql-postgres",
		dsn,
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
