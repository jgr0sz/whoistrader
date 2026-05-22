package main

import (
	"log"
	"os"

	"github.com/jgr0sz/whoistrader/endpoints"
	"github.com/joho/godotenv"
)

func main() {
	//Random IDs as an example (for now)
	steamIDs := []uint64{
		76561198332541485, //My ID
		76561199250221928, //ID of a user that has reversed, private CSFloat stall
	}

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, relying on environment variables")
	}

	//Registered endpoints
	registry := NewRegistry()
	registry.Register(&endpoints.CSFloatEndpoint{APIKey: os.Getenv("CSFLOAT_API_KEY")})
	registry.Register(&endpoints.ReverseWatchEndpoint{})

	//Looping through target IDs, forming aggregated profiles for each.
	for _, id := range steamIDs {
		if err := CreateProfile(id, registry); err != nil {
			log.Printf("%v", err)
		}
	}
}
