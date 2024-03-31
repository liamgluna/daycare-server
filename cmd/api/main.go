package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/liamgluna/daycare-server/internal/data"
	_ "github.com/lib/pq"
)

type config struct {
	port string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	allowCORS string
	jwtSecret string
}

type application struct {
	cfg    config
	logger *slog.Logger
	models data.Models
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
		os.Exit(1)
	}

	var cfg config
	flag.StringVar(&cfg.port, "port", os.Getenv("PORT"), "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DAYCARE_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.StringVar(&cfg.allowCORS, "allowCORS", os.Getenv("ALLOW_CORS"), "Allow CORS")
	flag.StringVar(&cfg.jwtSecret, "jwt-secret", os.Getenv("JWT_SECRET"), "JWT secret key")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database connection pool established")

	defer db.Close()

	app := &application{
		cfg:    cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

// openDB opens a new database connection pool using the configuration settings
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
