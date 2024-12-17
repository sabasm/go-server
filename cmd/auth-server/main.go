package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/sabasm/go-server/cmd/auth-server/service"
    "github.com/sabasm/go-server/internal/config"
    "github.com/sabasm/go-server/internal/logger"
    "github.com/sabasm/go-server/internal/server"
    "go.uber.org/zap"
)

func handleServerShutdown(srv server.ServerInterface, logger *zap.Logger) {
    logger.Info("Initiating graceful shutdown")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Error("Forced shutdown required", zap.Error(err))
    }
}

func main() {
    appConfig := config.LoadFromEnv()
    if err := appConfig.Validate(); err != nil {
        log.Fatalf("Configuration validation failed: %v", err)
    }

    logger, err := logger.NewLogger()
    if err != nil {
        log.Fatalf("Logger initialization failed: %v", err)
    }
    defer func() { _ = logger.Sync() }()

    srvCfg := &server.Config{
        Host:     appConfig.GetAppHost(),
        Port:     appConfig.GetAppPort(),
        BasePath: "/auth",
        Options: server.Options{
            ReadTimeout:    15 * time.Second,
            WriteTimeout:   15 * time.Second,
            IdleTimeout:    60 * time.Second,
            MaxHeaderBytes: 1 << 20,
        },
    }

    authService := service.NewAuthService(logger)
    authService.ConfigureRoutes(os.Getenv("M2M_SECRET_KEY"))

    srv := server.NewBuilder(srvCfg).
        WithLogger(logger).
        WithTimeout(
            srvCfg.Options.ReadTimeout,
            srvCfg.Options.WriteTimeout,
            srvCfg.Options.IdleTimeout,
        ).
        WithRoute("/auth", authService.GetHandler().ServeHTTP).
        Build()

    serverError := make(chan error, 1)
    go func() {
        logger.Info("Starting auth server",
            zap.String("host", srvCfg.Host),
            zap.Int("port", srvCfg.Port))
        if err := srv.Start(); err != nil && err != http.ErrServerClosed {
            logger.Error("Server error occurred", zap.Error(err))
            serverError <- err
        }
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    select {
    case <-sigChan:
        handleServerShutdown(srv, logger)
    case err := <-serverError:
        logger.Error("Fatal server error", zap.Error(err))
    }
}


