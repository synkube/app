package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/synkube/app/blueprint/data"
	coreData "github.com/synkube/app/core/data"
	"github.com/synkube/app/core/ginhelper"
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
	r := ginhelper.New()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Printf("HTTP server is running on %s...\n", addr)
	log.Fatal(r.Run(addr))
}

func startGraphQLServer(cfg coreData.ServerConfig, dm *data.DataModel) {
	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	r := gin.New()

	h := handler.New(&handler.Config{
		Schema:   data.NewGraphQLSchema(dm),
		Pretty:   true,
		GraphiQL: true,
	})
	r.POST("/graphql", gin.WrapH(h))

	// GraphQL Playground handler
	playgroundHandler := playground.Handler("GraphQL Playground", "/graphql")
	r.GET("/graphql", gin.WrapH(playgroundHandler))
	r.GET("/gq", gin.WrapH(playgroundHandler))

	log.Printf("GraphQL server is running on %s...\n", addr)
	log.Fatal(r.Run(addr))
}
