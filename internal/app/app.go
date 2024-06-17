package app

import (
    "log/slog"
    "time"

    "github.com/solloball/sso/internal/app/grpc"
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
    // TODO:: init storage

    //TODO:: init auth server

    grpcApp := grpcapp.New(log, grpcPort)

    return &App {
        GRPCApp: grpcApp,
    }
}
