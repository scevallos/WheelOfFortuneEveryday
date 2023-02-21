package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/scevallos/WheelOfFortuneEveryday/pkg/logging"
	"github.com/scevallos/WheelOfFortuneEveryday/pkg/roku"
	"github.com/scevallos/WheelOfFortuneEveryday/pkg/web"
)

const (
	port = ":8787"
)

var healthy int32

func main() {
	logger := logging.NewLogger()
	client, err := roku.NewClient(&roku.ClientOptions{
		Logger: logger,
		Config: &roku.Config{},
	})
	panicIfErr(err)

	service, err := web.NewService(client, logger, &healthy)
	panicIfErr(err)

	logger.Println("Starting HTTP server")
	server := &http.Server{
		Addr:         port,
		Handler:      service.GetRouter(),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at " + port)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", port, err)
	}

	<-done
	logger.Println("Server stopped")
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
