package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	"TMS.netjonin.net/internal/data"
	"TMS.netjonin.net/internal/jsonlog"
	"TMS.netjonin.net/internal/mailer"
	"TMS.netjonin.net/internal/vcs"
	
	_ "github.com/lib/pq"
)

var version = vcs.Version()

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}

	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}

	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}

	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {

	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port0")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	//flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://netjonin:pa55word@localhost/databasename", "PostgreSQL DSN")
	//flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("TMS_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags into the config struct.
	// Notice the default values that we're using?
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	// Create command line flags to read the setting values into the config struct.
	// Notice that we use true as the default for the 'enabled' setting?
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "7950205b99f449", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "a1d2481b1bc422", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "TMS <no-reply@TMS.netjonin.net>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

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

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
