package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	IsProd        bool
	HistoryLength int
}

func envPath() string {
	if len(os.Args) < 2 {
		return ".env"
	}
	return os.Args[1]
}

func LoadConfig() *Config {
	_ = godotenv.Load(envPath())
	isProd := os.Getenv("IS_PROD")
	fmt.Println("isProd: ", isProd)

	return &Config{
		IsProd:        os.Getenv("IS_PROD") == "true",
		HistoryLength: 100,
	}
}
