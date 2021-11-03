package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/eduardohoraciosanto/simple-blockchain/config"
	"github.com/eduardohoraciosanto/simple-blockchain/pkg/blockchain"
	"github.com/eduardohoraciosanto/simple-blockchain/pkg/health"
	"github.com/eduardohoraciosanto/simple-blockchain/transport"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := logrus.New()

	svc := blockchain.NewBlockChain(
		log.WithField("owner", "blockchain-service").Logger,
		config.GetVersion(),
	)

	hsvc := health.NewService(
		log.WithField("owner", "health-service").Logger,
	)

	httpTransportRouter := transport.NewHTTPRouter(svc, hsvc)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", config.GetPort()),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      httpTransportRouter,
	}
	log.WithField(
		"transport", "http").
		WithField(
			"port", config.GetPort()).
		Log(logrus.InfoLevel, "Transport Start")
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.WithField(
				"transport", "http").
				WithError(err).
				Log(logrus.ErrorLevel, "Transport Stopped")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Log(logrus.InfoLevel, "Service gracefully shutted down")
	os.Exit(0)
}
