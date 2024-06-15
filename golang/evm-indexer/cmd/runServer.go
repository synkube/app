package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	coreData "github.com/synkube/app/core/data"
	"github.com/synkube/app/core/ginhelper"
	"github.com/synkube/app/evm-indexer/data"
	"github.com/synkube/app/evm-indexer/graphql/graph"
)

func StartServers(servers []coreData.ServerConfig, bds *data.BlockchainDataStore) {
	for _, server := range servers {
		switch server.Type {
		case "http":
			go startHTTPServer(server)
		case "grpc":
			// Add gRPC server initialization here
		case "graphql":
			go startGraphQLServer(server, bds)
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
	r := ginhelper.New([]string{ginhelper.HealthCheckRoute, ginhelper.RobotsTxtRoute})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Printf("HTTP server is running on %s...\n", addr)
	log.Fatal(r.Run(addr))
}

func startGraphQLServer(cfg coreData.ServerConfig, bds *data.BlockchainDataStore) {
	addr := fmt.Sprintf(":%d", cfg.Port)
	r := ginhelper.New([]string{})

	// GraphQL handler
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{BDS: bds}}))

	// GraphQL Playground handler
	playgroundHandler := playground.Handler("GraphQL Playground", "/query")

	// Define routes
	r.GET("/", gin.WrapH(playgroundHandler))
	r.GET("/graphql", gin.WrapH(playgroundHandler))
	r.GET("/gq", gin.WrapH(playgroundHandler))
	r.POST("/query", gin.WrapH(srv))

	log.Printf("GraphQL server is running on %s...\n", addr)
	log.Fatal(r.Run(addr))
}
