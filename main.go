package main

import (
	"context"
	"log"
	"os"

	"github.com/jgr0sz/whoistrader/endpoints"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, relying on environment variables")
	}

	//Registered endpoints
	registry := NewRegistry()
	registry.Register(&endpoints.CSFloatEndpoint{})
	registry.Register(&endpoints.ReverseWatchEndpoint{})
	registry.Register(&endpoints.SteamInfoEndpoints{APIKey: os.Getenv("STEAM_API_KEY")})

	cmd := buildCLI(registry)
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
