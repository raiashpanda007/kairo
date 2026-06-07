package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type GRPCServer struct {
	Addr string
}

type Config struct {
	Server       GRPCServer
	IsProduction bool
}

func mustEnv(key string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}

func boolEnv(key string) bool {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	return val == "true" || val == "1"
}

func MustLoad() *Config {
	log.Println("loading config...")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	return &Config{
		Server: GRPCServer{
			Addr: mustEnv("GRPC_SERVER_ADDR"),
		},
		IsProduction: boolEnv("PRODUCTION"),
	}
}
