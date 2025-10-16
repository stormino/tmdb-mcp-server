package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	mode := flag.String("mode", getEnvOrDefault("MCP_MODE", "stdio"), "Transport mode: stdio or http")
	port := flag.String("port", getEnvOrDefault("HTTP_PORT", "8080"), "HTTP port (only used in http mode)")
	flag.Parse()

	tmdbAPIKey := os.Getenv("TMDB_API_KEY")
	if tmdbAPIKey == "" {
		log.Fatal("TMDB_API_KEY environment variable is required")
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "tmdb-api",
		Version: "v1.0.0",
	}, nil)

	tmdbServer := NewTMDBServer(tmdbAPIKey)
	tmdbServer.RegisterTools(server)

	switch *mode {
	case "stdio":
		log.Println("Starting TMDB MCP server in stdio mode (local)")
		transport := &mcp.StdioTransport{}
		log.Println("TMDB MCP Server starting...")
		if err := server.Run(context.Background(), transport); err != nil {
			log.Fatal(err)
		}

	case "http":
		log.Printf("Starting TMDB MCP server in HTTP mode (remote) on port %s\n", *port)

		handler := mcp.NewSSEHandler(func(req *http.Request) *mcp.Server {
			return server
		}, nil)

		address := ":" + *port
		log.Printf("TMDB MCP Server starting on http://localhost%s\n", address)
		if err := http.ListenAndServe(address, handler); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("Unknown mode: %s. Use 'stdio' or 'http'", *mode)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
