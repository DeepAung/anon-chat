package config

import (
	"log"
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
	path := envPath()
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatal("load dotenv failed: ", err)
	}

	err = godotenv.Load(path)
	if err != nil {
		log.Fatal("load dotenv failed: ", err)
	}

	return &Config{
		IsProd:        envMap["isProd"] == "true",
		HistoryLength: 100,
	}
}
