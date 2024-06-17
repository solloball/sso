package config

import (
    "time"
    "os"

    "github.com/ilyakaznacheev/cleanenv"
    "github.com/joho/godotenv"
)

type Config struct {
    Env string `yaml:"env" env-default:"prod"`
    StoragePath string `yaml:"storage_path" env-required:"true"`
    TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
    GRPC GRPCConfig `yaml:"grpc" env-required:"true"`
}

type GRPCConfig struct {
    Port int `yaml:"port"`
    Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
    path := fetchConfigPath()

    if _, err := os.Stat(path); os.IsNotExist(err) {
        panic("config file does not exist " + path) 
    }
    
    var cfg Config

    if err := cleanenv.ReadConfig(path, &cfg); err != nil {
        panic("failed to read config" + err.Error())
    }

    return &cfg
}

func fetchConfigPath() string {
    if err := godotenv.Load(".env"); err != nil {
        panic(err)
    }

    return os.Getenv("CONFIG_PATH")
}
