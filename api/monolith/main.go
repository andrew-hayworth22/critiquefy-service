package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/andrew-hayworth22/critiquefy-service/api/monolith/debug"
	"github.com/andrew-hayworth22/critiquefy-service/api/monolith/mux"
	"github.com/andrew-hayworth22/critiquefy-service/app/auth"
	"github.com/andrew-hayworth22/critiquefy-service/business/data/sqldb"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/keystore"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/logger"
	"github.com/andrew-hayworth22/critiquefy-service/foundation/web"
	"github.com/ardanlabs/conf/v3"
	"github.com/joho/godotenv"
)

func main() {
	// -----------------------------------------------------------------
	// Configure logger

	events := logger.Events{}
	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}
	log := logger.NewWithEvents(os.Stdout, logger.LevelInfo, "CRITIQUEFY", traceIDFn, events)

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -----------------------------------------------------------------
	// Configuration

	err := godotenv.Load()
	if err != nil {
		log.Error(ctx, "startup", err)
	}

	cfg := struct {
		Version struct {
			Build       string `conf:"default:'DEV'"`
			Number      string `conf:"default:'0.0.1'"`
			Description string `conf:"default:'Critiquefy DEV 0.0.1'"`
		}
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:3010"`
			CORSAllowedOrigins []string      `conf:"default:*,mask"`
		}
		Auth struct {
			KeysFolder string `conf:"default:zarf/keys/"`
			ActiveKID  string `conf:"default:1eba8606-9e70-411c-9d2f-e922431cfa37"`
			Issuer     string `conf:"default:critiquefy"`
		}
		DB struct {
			URL string `conf:"default:'postgresql://user:password@localhost/critiquefy'"`
		}
	}{}

	const prefix = "CRITIQUEFY"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -----------------------------------------------------------------
	// Starting App

	log.Info(ctx, "starting service", "status", "initializing authentication")

	ks := keystore.New()
	if _, err := ks.LoadByFileSystem(os.DirFS(cfg.Auth.KeysFolder)); err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config{
		Log:       log,
		KeyLookup: ks,
	}

	auth, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// -----------------------------------------------------------------
	// Database Support

	log.Info(ctx, "startup", "status", "initializing database support")

	db, err := sqldb.Open(ctx, sqldb.Config{
		URL: cfg.DB.URL,
	})
	if err != nil {
		return fmt.Errorf("connecting to DB: %w", err)
	}
	defer db.Close()

	// -----------------------------------------------------------------
	// Starting Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost)
		}
	}()

	// -----------------------------------------------------------------
	// Starting API Service

	log.Info(ctx, "startup", "status", "initializing API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.WebAPI(cfg.Version.Build, log, db, auth, shutdown),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -----------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
