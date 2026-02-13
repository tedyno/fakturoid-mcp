package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/fakturoid-mcp/config"
	"github.com/tedyno/fakturoid-mcp/fakturoid"
	"github.com/tedyno/fakturoid-mcp/tools"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	client := fakturoid.NewClient(cfg.ClientID, cfg.ClientSecret, cfg.Slug)

	s := server.NewMCPServer(
		"fakturoid-mcp",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	tools.RegisterAll(s, client)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
