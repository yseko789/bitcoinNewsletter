package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"github.com/yseko789/bitcoinSummarize/internal/data"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, os.Getenv("projectID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	logger.Info("firestore connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(client),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
