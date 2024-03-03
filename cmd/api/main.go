package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	port string
	db   struct {
		dsn string
	}
	allowCORS string
}

type application struct {
	cfg    config
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
		os.Exit(1)
	}

	var cfg config
	flag.StringVar(&cfg.port, "port", os.Getenv("PORT"), "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("SCHOLARLY_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.allowCORS, "allowCORS", os.Getenv("ALLOW_CORS"), "Allow CORS")


	flag.Parse()

	app := &application{
		cfg:    cfg,
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
