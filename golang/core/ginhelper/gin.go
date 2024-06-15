package ginhelper

import (
	"bytes"
	_ "embed"
	"net/http"
	"time"

	"github.com/ipfs/go-log"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

//go:embed robots.txt
var robots []byte

const (
	HealthCheck     = "/health"
	RobotsTxt       = "/robots.txt"
	RequestIDHeader = "X-Request-ID"
)

var HealthCheckRoute = HealthCheck
var RobotsTxtRoute = RobotsTxt

var bootTime = time.Now()

var logger = log.Logger("gin")

// New creates a new gin server with some sensible defaults.
func New(routes []string) *gin.Engine {
	server := gin.New()
	server.ContextWithFallback = true

	configureMiddleware(server)
	configureRoutes(server, routes)

	return server
}

// configureMiddleware sets up the middleware for the Gin server.
func configureMiddleware(server *gin.Engine) {
	server.Use(helmet.Default())
	server.Use(gin.Recovery())
	server.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"*"},
		AllowMethods:    []string{"GET", "PUT", "POST", "PATCH", "DELETE", "OPTIONS"},
		MaxAge:          12 * time.Hour,
	}))

	server.Use(requestid.New(
		requestid.WithCustomHeaderStrKey(RequestIDHeader),
		requestid.WithGenerator(func() string {
			return uuid.New().String()
		})))

	server.Use(ginzap.GinzapWithConfig(logger.Desugar(), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		Context: func(c *gin.Context) (fields []zapcore.Field) {
			requestID := c.GetHeader(RequestIDHeader)
			fields = append(fields, zapcore.Field{
				Key:    "request-id",
				Type:   zapcore.StringType,
				String: requestID,
			})

			return fields
		},
	}))

	server.Use(setRequestIDHeader)
	server.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))
}

// setRequestIDHeader sets the request ID header if the client didn't provide one.
func setRequestIDHeader(c *gin.Context) {
	if c.Request.Header.Get(RequestIDHeader) == "" {
		c.Request.Header.Set(RequestIDHeader, c.Writer.Header().Get(RequestIDHeader))
	}
}

// configureRoutes sets up the routes for the Gin server.
func configureRoutes(server *gin.Engine, routes []string) {
	for _, route := range routes {
		if route == HealthCheck {
			server.GET(HealthCheck, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "UP"})
			})
		} else if route == RobotsTxt {
			server.GET(RobotsTxt, func(c *gin.Context) {
				reader := bytes.NewReader(robots)
				c.Header(headers.ContentType, gin.MIMEPlain)
				http.ServeContent(c.Writer, c.Request, "robots.txt", bootTime, reader)
			})
		}
	}
}
