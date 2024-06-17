package main

import (
    "os"
    "log/slog"

    "github.com/solloball/sso/internal/config"
    "github.com/solloball/sso/internal/app"
)

const (
    envLocal = "local"
    envDev = "dev"
    envProd = "prod"
)

func main() {
    cfg := config.MustLoad()


    log := setupLogger(cfg.Env)

    log.Info("starting application", slog.String("env", cfg.Env))

    application := app.New(
        log,
        cfg.GRPC.Port,
        cfg.StoragePath,
        cfg.TokenTTL,
    )

    application.GRPCApp.MustRun()

    //TODO:: run grpc service
}

func setupLogger(env string) *slog.Logger {
    var log *slog.Logger

    switch env {
    case envLocal:
        log = slog.New(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    case envDev:
        log = slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    case envProd:
        log = slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
        )
    }

    return log
}
