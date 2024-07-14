package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sherpaurgen/argus/internal/data"
	"github.com/sherpaurgen/argus/internal/jsonlog"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct { //data source name (DSN)
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxLifeTime  int
	}
	limiter struct {
		rps     float64 //request per sec
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:password12345@localhost/pgargusdb?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.IntVar(&cfg.db.maxLifeTime, "db-max-life-time", 10, "PostgreSQL max connection life time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	//initialize logger, the prefix is empty hence its just empty quote ""
	//logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := OpenDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.recoverPanic(app.rateLimit(app.routes())),
		IdleTimeout:  time.Minute,
		ErrorLog:     log.New(logger, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP
	logger.PrintInfo("starting server", map[string]string{"addr": srv.Addr, "env": cfg.env})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, map[string]string{"error": err.Error()})

}

func OpenDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	durationStr := fmt.Sprintf("%dm", cfg.db.maxLifeTime) // "5m" for 5 minutes
	duration, _ := time.ParseDuration(durationStr)
	db.SetConnMaxLifetime(duration * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
