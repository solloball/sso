package main

import (
    "os"
    "log/slog"

    "github.com/solloball/sso/internal/config"
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
    //TODO:: init logger

    //TODO:: init app

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
