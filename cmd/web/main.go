package main

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/patrickarmengol/coffeetanuki/internal/data"
	"github.com/patrickarmengol/coffeetanuki/internal/vcs"
)

var version = vcs.Version()

type config struct {
	env    string // dev, staging, prod
	server struct {
		port         int
		idleTimeout  time.Duration
		readTimeout  time.Duration
		writeTimeout time.Duration
	}
	db struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type application struct {
	config        config
	formDecoder   *form.Decoder
	logger        *slog.Logger
	models        data.Models
	templateCache map[string]*template.Template
	wg            sync.WaitGroup
}

func main() {
	// parse commandline flags
	var cfg config
	flag.StringVar(&cfg.env, "env", "development", "environment (development|staging|production)")

	flag.IntVar(&cfg.server.port, "server-port", 4000, "server port")
	flag.DurationVar(&cfg.server.idleTimeout, "server-idle-timeout", time.Minute, "server idle timeout")
	flag.DurationVar(&cfg.server.readTimeout, "server-read-timeout", 5*time.Second, "server read timeout")
	flag.DurationVar(&cfg.server.writeTimeout, "server-write-timeout", 10*time.Second, "server write timeout")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open conections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle conections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	displayVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	// if version flag, print version and exit
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
	}

	// initialize structured lgr; writes to stdout
	lgr := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// initialize db conn pool
	db, err := openDB(cfg.db)
	if err != nil {
		lgr.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	lgr.Info("database connection pool established")

	// initialize models by providing db conn pool
	mdls := data.NewModels(db)

	// initialize template cache
	tmpls, err := newTemplateCache()
	if err != nil {
		lgr.Error(err.Error())
		os.Exit(1)
	}

	// initialize form decoder
	fdcdr := form.NewDecoder()

	// construct application
	app := &application{
		config:        cfg,
		formDecoder:   fdcdr,
		logger:        lgr,
		models:        mdls,
		templateCache: tmpls,
	}

	// start service
	err = app.serve()

	// on service end
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
