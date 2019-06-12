package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"

	"github.com/michaljemala/payments-sample/internal/doc"
	"github.com/michaljemala/payments-sample/internal/migrate/postgres"
	"github.com/michaljemala/payments-sample/pkg/payments"
)

var (
	flagAddr         = flag.String("http", "localhost:8080", "API server address")
	flagDsn          = flag.String("database", "", "Database server connect string")
	flagMigrationDir = flag.String("migrations", "", "Location of the migration files")
	flagDocs         = flag.Bool("docs", true, "")
)

func main() {
	flag.Parse()

	logger := initLogging()

	if *flagMigrationDir != "" {
		m := postgres.NewMigration(postgres.Config{
			DSN:          *flagDsn,
			DatabaseName: "payments",
			SourceURL:    *flagMigrationDir,
			Logger:       logger,
		})
		err := m.Do()
		if err != nil {
			logger.Fatalf("unable to migrate database: %v", err)
		}
	}

	api, close := initPaymentAPI(logger)
	defer close()

	router := initRouter(api.Prefix(), api, logger, *flagDocs)

	initAndStartServer(router, logger)
}

func initLogging() *log.Logger {
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC
	log.SetFlags(flags)
	return log.New(os.Stderr, "", flags)
}

func initPaymentAPI(logger *log.Logger) (*payments.API, func()) {
	api, err := payments.NewAPI(payments.Config{
		Prefix: "/",
		Driver: "postgres",
		DSN:    *flagDsn,
		Logger: logger,
	})
	if err != nil {
		logger.Fatalf("unable to initialize payment API: %v", err)
	}
	return api, func() {
		err := api.Close()
		if err != nil {
			logger.Printf("unable to gracefully close payment API: %v", err)
		}
	}
}

func initRouter(pattern string, handler http.Handler, logger *log.Logger, withDocs bool) http.Handler {
	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger}),
		middleware.Recoverer,
	)
	router.Mount(pattern, handler)
	if withDocs {
		router.Mount(doc.Prefix, doc.Handler())
	}

	return router
}

func initAndStartServer(router http.Handler, logger *log.Logger) {
	server := &http.Server{
		Addr:    *flagAddr,
		Handler: router,
	}

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		logger.Printf("server shutting down: signal received: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			logger.Printf("server shutdown failed: %v", err)
		}
	}()

	logger.Printf("server starting up: %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			logger.Fatalf("server startup failed: %v", err)
		}
	}
}
