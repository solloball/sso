package app

import (
    "log/slog"
    "time"

    "github.com/solloball/sso/internal/app/grpc"
    "github.com/solloball/sso/internal/storage/sqlite"
    "github.com/solloball/sso/internal/services/auth"
)

type App struct {
    GRPCApp * grpcapp.App
}

func New(
     log *slog.Logger,
     grpcPort int,
     storagePath string,
     tokenTTL time.Duration,
) *App {
    storage, err := sqlite.New(storagePath)
    if err != nil {
        panic(err)
    }

    authService := auth.New(log, storage, storage, storage, tokenTTL)

    grpcApp := grpcapp.New(log, authService, grpcPort)

    return &App {
        GRPCApp: grpcApp,
    }
}
