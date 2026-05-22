package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/jgr0sz/whoistrader/core"
)

//Wrapper struct to consolidate our endpoints.
type Registry struct {
	endpoints []core.Endpoint
}

//Constructor for our Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

//Adds a new Endpoint to the registry.
func (r *Registry) Register(e core.Endpoint) {
	r.endpoints = append(r.endpoints, e)
}

//Main aggregator logic that concurrently extracts Endpoints and their functions, retrieving responses/errors.
//Mutexes are used here in order to ensure race conditions between goroutines accessing a shared struct don't occur.
func AggregateTraderProfile(steamID uint64, registry *Registry) (*core.AggregatedTraderProfile, error) {
	profile := &core.AggregatedTraderProfile{
		SteamID: strconv.FormatUint(steamID, 10),
		Data:    make(map[string]any),
	}

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, endpoint := range registry.endpoints {
		wg.Go(func() {
			result, err := endpoint.Fetch(steamID)
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				profile.Errors = append(profile.Errors, fmt.Sprintf("%s: %v", endpoint.Name(), err))
				return
			}
			profile.Data[endpoint.Name()] = result
		})
	}
	wg.Wait()

	if len(profile.Data) == 0 {
		return nil, fmt.Errorf("All APIs failed to fetch a response: %v", profile.Errors)
	}
	return profile, nil
}

//Wrapper function to invoke the aggregator, handle errors, and output JSON aggregator response/s.
func CreateProfile(steamID uint64, registry *Registry) error {
	profile, err := AggregateTraderProfile(steamID, registry)
	if err != nil {
		return fmt.Errorf("Aggregation failed for steamID %d: %v\n", steamID, err)
	}

	if len(profile.Errors) > 0 {
		log.Printf("Some sources failed for steamID %d: %v", steamID, err)
	}
	return json.NewEncoder(os.Stdout).Encode(profile)
}
