package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/jgr0sz/whoistrader/core"
)

type Registry struct {
	endpoints []core.Endpoint
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(e core.Endpoint) {
	r.endpoints = append(r.endpoints, e)
}

//We extract data from our registered endpoints and their fetch functions into the user profile.
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
