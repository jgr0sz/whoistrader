package core

//Interface for API endpoints.
type Endpoint interface {
	Fetch(steamID uint64) (any, error)
	Name() string
}

//Aggregator struct for APIs.
type AggregatedTraderProfile struct {
	SteamID string `json:"steam_id"`
	Data map[string]any `json:"data"`
	Errors []string `json:"errors,omitempty"`
}