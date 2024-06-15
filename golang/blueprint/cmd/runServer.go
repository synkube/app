package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/graphql-go/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/synkube/app/blueprint/data"
	coreData "github.com/synkube/app/core/data"
)

func StartServers(servers []coreData.ServerConfig, dm *data.DataModel) {
	for _, server := range servers {
		switch server.Type {
		case "http":
			go startHTTPServer(server)
		case "grpc":
			// Add gRPC server initialization here
		case "graphql":
			go startGraphQLServer(server, dm)
		case "websocket":
			// Add WebSocket server initialization here
		default:
			log.Printf("Unsupported server type: %s", server.Type)
		}
	}

	// Keep the main routine running
	select {}
}

func startHTTPServer(cfg coreData.ServerConfig) {
	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("HTTP server is running on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func startGraphQLServer(cfg coreData.ServerConfig, dm *data.DataModel) {
	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	h := handler.New(&handler.Config{
		Schema:   data.NewGraphQLSchema(dm),
		Pretty:   true,
		GraphiQL: true, // Enable GraphiQL interface
	})
	http.Handle("/graphql", h)

	playgroundHandler := playground.Handler("GraphQL Playground", "/graphql")
	http.Handle("/gqi", playgroundHandler)

	log.Printf("GraphQL server is running on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
