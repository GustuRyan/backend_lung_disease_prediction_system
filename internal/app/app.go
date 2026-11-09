package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/config"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/database"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/handler"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/repository"
	"github.com/GustuRyan/backend_lung_disease_prediction_system/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	db, err := database.Open(cfg.DatabaseDSN, cfg.Environment)
	if err != nil {
		return err
	}
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery(), requestID(), requestLogger(logger), securityHeaders())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	handler.New(service.New(repository.New(db))).RegisterRoutes(router)
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access database pool: %w", err)
	}
	router.GET("/health/ready", func(c *gin.Context) {
		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("http server started", "address", cfg.Address, "environment", cfg.Environment)
		serverErrors <- server.ListenAndServe()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case signal := <-shutdownSignal:
		logger.Info("shutdown signal received", "signal", signal.String())
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("serve http: %w", err)
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	return sqlDB.Close()
}

func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()
		attributes := []any{
			"request_id", c.GetString("request_id"),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(started).Milliseconds(),
			"client_ip", c.ClientIP(),
		}
		if len(c.Errors) > 0 {
			attributes = append(attributes, "error", c.Errors.String())
			logger.Error("http request failed", attributes...)
			return
		}
		logger.Info("http request", attributes...)
	}
}

func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" || len(id) > 128 {
			bytes := make([]byte, 16)
			if _, err := rand.Read(bytes); err == nil {
				id = hex.EncodeToString(bytes)
			} else {
				id = fmt.Sprintf("fallback-%d", time.Now().UnixNano())
			}
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Next()
	}
}
