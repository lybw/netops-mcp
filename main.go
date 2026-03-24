package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lybw/netops-mcp/tools"
	"github.com/mark3labs/mcp-go/server"
)

const version = "0.1.0"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("netops-mcp %s\n", version)
		return
	}

	s := server.NewMCPServer(
		"netops-mcp",
		version,
		server.WithToolCapabilities(true),
		server.WithRecovery(),
		server.WithInstructions("Network operations MCP server. Provides tools for LAN device discovery, port scanning, OS fingerprinting, and MAC vendor lookup."),
	)

	tools.RegisterAll(s)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
