package main

import (
	"os"

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

	// Run the server
	srv.Run()

}
