package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// setup server with configuration; handle shutdown gracefully
func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.server.port),
		Handler:      app.routes(),
		IdleTimeout:  app.config.server.idleTimeout,
		ReadTimeout:  app.config.server.readTimeout,
		WriteTimeout: app.config.server.writeTimeout,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	// define chan to receive errors returned by Shutdown()
	shutdownError := make(chan error)

	// goroutine to listen for signals
	go func() {
		// define new buffered signal chan
		quit := make(chan os.Signal, 1)

		// listen for incoming SIGINT, SIGTERM and pass to quit chan
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// block until quit chan is read
		s := <-quit

		// log received signal
		app.logger.Info("shutting down server", "signal", s.String())

		// create context with 30s timeout
		// grace period for existing requests to complete
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// wait for any background tasks to finish
		app.logger.Info("completing background tasks", "addr", srv.Addr)
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	// expect http.ErrServerClosed for graceful shutdown start when Shutdown() called
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// return val from Shutdown() indicates success or not
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
