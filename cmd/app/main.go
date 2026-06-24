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

router := gin.New()
router.Use(gin.Recovery())
router.Use(requestLogger(log))

router.GET("/health", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{
"status":  "ok",
"service": "subscriptions-api",
})
})

server := &http.Server{
Addr:         ":" + cfg.AppPort,
Handler:      router,
ReadTimeout:  10 * time.Second,
WriteTimeout: 10 * time.Second,
}

ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
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

shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
