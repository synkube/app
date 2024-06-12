package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/synkube/app/core/data"
)

func StartServers(servers []data.ServerConfig) {
	for _, server := range servers {
		switch server.Type {
		case "http":
			go startHTTPServer(server)
		case "grpc":
			// Add gRPC server initialization here
		case "graphql":
			// Add GraphQL server initialization here
		case "websocket":
			// Add WebSocket server initialization here
		default:
			log.Printf("Unsupported server type: %s", server.Type)
		}
	}

	// Keep the main routine running
	select {}
}

func startHTTPServer(cfg data.ServerConfig) {
	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("HTTP server is running on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
