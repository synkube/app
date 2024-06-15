package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/synkube/app/core/data"
	"github.com/synkube/app/core/ginhelper"
)

func StartHTTPServer(cfg data.ServerConfig) {
	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("HTTP server is running on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func StartHTTPGinServer(cfg data.ServerConfig) {
	addr := fmt.Sprintf("localhost:%d", cfg.Port)
	r := ginhelper.New([]string{ginhelper.HealthCheckRoute, ginhelper.RobotsTxtRoute})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Printf("HTTP server is running on %s...\n", addr)
	log.Fatal(r.Run(addr))
}
