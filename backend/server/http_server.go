package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Simaky/go-github-tracker/backend/consts"
	"github.com/Simaky/go-github-tracker/backend/server/handlers"
)

const (
	shutdownTimeout   = 30 * time.Second
	readHeaderTimeout = 10 * time.Second
	maxHeaderBytes    = 1 << 16
)

func (s *server) runHTTP(version string) error {
	h := handlers.New(s.app)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(RequestID(), RequestLogger(), gin.Recovery())

	router.GET("/uptime", h.Uptime)
	router.GET("/version", func(c *gin.Context) { c.String(http.StatusOK, version) })
	router.Any("/debug/pprof/*profile", gin.WrapH(http.DefaultServeMux)) // pprof self-registers via main's blank import

	// Resource API, mounted under /api per the assignment.
	api := router.Group("/api")
	{
		api.POST("/repos", h.CreateRepo)
		api.GET("/repos", h.ListRepos)
		api.GET("/repos/:id", h.GetRepo)
		api.GET("/metrics", h.TotalMetrics)
		api.PATCH("/repos/:id", h.UpdateNotes)
		api.POST("/repos/:id/refresh", h.RefreshRepo)
		api.DELETE("/repos/:id", h.DeleteRepo)
	}

	httpSrv := &http.Server{
		Addr:              s.addressOf(),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
	return runWithGracefulShutdown(httpSrv)
}

func runWithGracefulShutdown(httpSrv *http.Server) error {
	idle := make(chan error, 1)
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		idle <- httpSrv.Shutdown(ctx)
	}()

	log.Printf("%s listening on http://%s/", consts.ServiceName, httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http listen: %w", err)
	}
	return <-idle
}
