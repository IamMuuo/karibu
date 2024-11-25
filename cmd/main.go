package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/iammuuo/karibu/config"
	"github.com/iammuuo/karibu/server"
)

func main() {

	// Load Configurations
	cfg, err := config.LoadDefaultConfigs()

	if err != nil {
		log.Errorf("Failed to load configuration file with error: %v", err)
		os.Exit(2)
	}

	// Start a new server
	var srv *server.Server
	srv, err = server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create new server with error: %v\n", err)
	}

	// create a channel for interrupt handling
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run the server
	go srv.Run()

	// Here we wait for the done signal: this can be either an interrupt, or
	// the server shutting down for any other reason.
	<-done

	// When it arrives, we create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() { cancel() }()

	// When we start the shutdown, the server will no longer accept new
	// connections, but will wait as much as the given context allows for the
	// active connections to finish.
	// After the timeout, it shuts down anyway.
	srv.Shutdown(ctx)

}
