package main

import (
    "fmt"

    "github.com/solloball/sso/internal/config"
)

func main() {
    cfg := config.MustLoad()

    //TODO:: remove it
    fmt.Println(cfg)

    //TODO:: init logger

    //TODO:: init app

    //TODO:: run grpc service
}
