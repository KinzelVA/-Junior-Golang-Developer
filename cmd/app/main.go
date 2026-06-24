package main

import (
"context"
"errors"
"log/slog"
"net/http"
"os"
"os/signal"
"syscall"
"time"

"github.com/gin-gonic/gin"

"github.com/KinzelVA/-Junior-Golang-Developer/internal/config"
"github.com/KinzelVA/-Junior-Golang-Developer/internal/db"
appLogger "github.com/KinzelVA/-Junior-Golang-Developer/internal/logger"
)

func main() {
cfg, err := config.Load()
if err != nil {
panic(err)
}

log := appLogger.New(cfg.AppEnv)
slog.SetDefault(log)

if cfg.AppEnv == "production" {
gin.SetMode(gin.ReleaseMode)
}

rootCtx := context.Background()

dbCtx, dbCancel := context.WithTimeout(rootCtx, 10*time.Second)
defer dbCancel()

postgresPool, err := db.NewPostgresPool(dbCtx, cfg.DatabaseURL())
if err != nil {
log.Error("failed to connect to PostgreSQL", slog.String("error", err.Error()))
os.Exit(1)
}
defer postgresPool.Close()

log.Info("connected to PostgreSQL")

router := gin.New()
router.Use(gin.Recovery())
router.Use(requestLogger(log))

router.GET("/health", func(c *gin.Context) {
ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
defer cancel()

if err := postgresPool.Ping(ctx); err != nil {
log.Error("database health check failed", slog.String("error", err.Error()))

c.JSON(http.StatusServiceUnavailable, gin.H{
"status":   "error",
"service":  "subscriptions-api",
"database": "unavailable",
})
return
}

c.JSON(http.StatusOK, gin.H{
"status":   "ok",
"service":  "subscriptions-api",
"database": "ok",
})
})

server := &http.Server{
Addr:         ":" + cfg.AppPort,
Handler:      router,
ReadTimeout:  10 * time.Second,
WriteTimeout: 10 * time.Second,
}

ctx, stop := signal.NotifyContext(rootCtx, os.Interrupt, syscall.SIGTERM)
defer stop()

go func() {
log.Info("starting HTTP server", slog.String("addr", server.Addr))

if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
log.Error("failed to start HTTP server", slog.String("error", err.Error()))
os.Exit(1)
}
}()

<-ctx.Done()

log.Info("shutting down HTTP server")

shutdownCtx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
defer cancel()

if err := server.Shutdown(shutdownCtx); err != nil {
log.Error("failed to shutdown HTTP server", slog.String("error", err.Error()))
os.Exit(1)
}

log.Info("HTTP server stopped")
}

func requestLogger(log *slog.Logger) gin.HandlerFunc {
return func(c *gin.Context) {
start := time.Now()

c.Next()

log.Info(
"http request",
slog.String("method", c.Request.Method),
slog.String("path", c.Request.URL.Path),
slog.Int("status", c.Writer.Status()),
slog.Duration("duration", time.Since(start)),
slog.String("client_ip", c.ClientIP()),
)
}
}
