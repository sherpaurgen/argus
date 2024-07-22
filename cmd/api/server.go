package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", app.config.port),
		Handler:     app.recoverPanic(app.rateLimit(app.routes())),
		IdleTimeout: time.Minute,
		// ErrorLog:     log.New(logger, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		//below incoming Signal Interrupt + signal terminate which gives a process the opportunity to perform clean-up tasks before it ends. those signal are relayed to channel named "quit"
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		//blocking code below
		//read quit channel
		s := <-quit

		app.logger.PrintInfo("Shutting down: caught signal", map[string]string{
			"signal": s.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		shutdownError <- srv.Shutdown(ctx)
		os.Exit(0)

	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})
	// Start the server as normal, returning any error.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil
}
