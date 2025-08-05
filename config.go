package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	GoogleSecret string
	GoogleKey    string
}

func loadEnv() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()
}

func getConfig() (config, error) {
	loadEnv()

	googleSecret := os.Getenv("GOOGLE_SECRET")
	if googleSecret == "" {
		return config{}, errors.New("GOOGLE_SECRET is not set")
	}
	googleKey := os.Getenv("GOOGLE_KEY")
	if googleKey == "" {
		return config{}, errors.New("GOOGLE_KEY is not set")
	}
	return config{
		GoogleSecret: googleSecret,
		GoogleKey:    googleKey,
	}, nil
}
