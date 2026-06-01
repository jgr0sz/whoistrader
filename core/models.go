package core

import "context"

//Interface for API endpoints.
type Endpoint interface {
	Fetch(ctx context.Context, steamID uint64) (any, error)
	Name() string
}

//Aggregator struct used to gather data across APIs for the aggregated profile. Data holds a generic value as it stores pointers to each type of API struct fetched.
type AggregatedTraderProfile struct {
	SteamID string         `json:"steam_id"`
	Data    map[string]any `json:"data"`
	Errors  []string       `json:"errors,omitempty"`
}
