package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/basic/docs"

	logger "github.com/charmbracelet/log"
)

var (
	Version   string
	GitCommit string
	BuildDate string
	log       = logger.NewWithOptions(os.Stderr, logger.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Level:           logger.DebugLevel,
	})
)

func RunServer(port int, version string) {
	addr := fmt.Sprintf(":%d", port)
	router := GinRouter(port)
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	Version = version
	log.Infof("Starting server on %s", addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down http server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Info("Server exiting")
}

// @contact.name   Ivan Pedrazas
// @contact.url    https://i2.alacasa.uk
// @contact.email  ipedrazas@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
func GinRouter(port int) *gin.Engine {
	docs.SwaggerInfo.Title = "Ivan's Internal Platform API"
	docs.SwaggerInfo.Description = "API to create, run and manage Applications."
	docs.SwaggerInfo.Version = "0.1.0"
	addr := fmt.Sprintf("localhost:%d", port)
	docs.SwaggerInfo.Host = addr
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router := gin.New()

	// Set a lower memory limit for multipart forms
	router.MaxMultipartMemory = 100 << 20 // 100 MB

	// Custom Logger
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s |%s %d %s| %s |%s %s %s %s | %s | %s | %s\n",
			param.TimeStamp.Format(time.RFC1123),
			param.StatusCodeColor(),
			param.StatusCode,
			param.ResetColor(),
			param.ClientIP,
			param.MethodColor(),
			param.Method,
			param.ResetColor(),
			param.Path,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Use(cors.Default())

	// Prometheus
	prom := NewPrometheus("gin")
	prom.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			if p.Key == `uid` {
				url = strings.Replace(url, p.Value, ":uid", 1)
			}
			if p.Key == `id` {
				url = strings.Replace(url, p.Value, ":id", 1)
			}
		}
		return url
	}
	metricUrl := fmt.Sprintf("http://localhost:%d/metrics", port)
	pushGWUrl := os.Getenv("PUSHGATEWAY_URL")
	pushIntervalStr := os.Getenv("PUSHGATEWAY_INTERVAL")
	pushInterval, err := time.ParseDuration(pushIntervalStr)
	if err != nil {
		pushInterval = time.Duration(180)
	}
	prom.SetPushGateway(pushGWUrl, metricUrl, pushInterval)
	prom.Use(router)

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	// router.Use(gin.Logger())
	AddRoutes(router)

	return router
}
