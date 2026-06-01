package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
func AggregateTraderProfile(ctx context.Context, steamID uint64, registry *Registry) (*core.AggregatedTraderProfile, error) {
	profile := &core.AggregatedTraderProfile{
		SteamID: strconv.FormatUint(steamID, 10),
		Data:    make(map[string]any),
	}

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, endpoint := range registry.endpoints {
		wg.Go(func() {
			result, err := endpoint.Fetch(ctx, steamID)
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
		return nil, fmt.Errorf("all APIs failed to fetch a response: %s", strings.Join(profile.Errors, "; "))
	}
	return profile, nil
}

//Wrapper function to invoke the aggregator, handle errors, and output JSON aggregator response/s.
func CreateProfile(ctx context.Context, steamID uint64, registry *Registry, savepath string) error {
	profile, err := AggregateTraderProfile(ctx, steamID, registry)
	if err != nil {
		return fmt.Errorf("Aggregation failed for steamID %d: %s\n", steamID, err)
	}
	if len(profile.Errors) > 0 {
		log.Printf("Some sources failed for steamID %d: %s\n\n", steamID, strings.Join(profile.Errors, "; "))
	}

	output, err := json.MarshalIndent(profile, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal profile for steamID %d: %s", steamID, err)
	}

	//Default: will return JSON response to standard output
	if savepath == "" {
		fmt.Fprintln(os.Stdout, string(output))
	} else {
		//Checks for existing savepath, makes file otherwise and writes to it
		file, err := os.OpenFile(savepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %s", savepath, err)
		}
		defer file.Close()

		if _, err := fmt.Fprintln(file, string(output)); err != nil {
			return fmt.Errorf("failed to write profile to %s: %s", savepath, err)
		}
	}
	return nil
}
