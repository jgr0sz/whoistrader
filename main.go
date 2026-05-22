package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/jgr0sz/whoistrader/endpoints"
	"github.com/joho/godotenv"
)

//Random ID as an example (for now)
const steamID uint64 = 76561198332541485

func main() {
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found, relying on environment variables")
    }

    registry := NewRegistry()
    registry.Register(&endpoints.CSFloatEndpoint{APIKey: os.Getenv("CSFLOAT_API_KEY")})

    profile, err := AggregateTraderProfile(steamID, registry)
    if err != nil {
        log.Fatalf("Aggregation failed: %v", err)
    }

    if len(profile.Errors) > 0 {
        log.Printf("Some sources failed: %v", profile.Errors)
    }

    if err := json.NewEncoder(os.Stdout).Encode(profile); err != nil {
        log.Fatalf("Failed to encode profile: %v", err)
    }
}