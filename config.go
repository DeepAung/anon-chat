package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	IsProd bool
}

func envPath() string {
	if len(os.Args) < 2 {
		return ".env"
	}
	return os.Args[1]
}

func LoadConfig() *config {
	path := envPath()
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatal("load dotenv failed: ", err)
	}

	err = godotenv.Load(path)
	if err != nil {
		log.Fatal("load dotenv failed: ", err)
	}

	return &config{
		IsProd: envMap["isProd"] == "true",
	}
}
