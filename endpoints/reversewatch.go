package endpoints

import (
	"encoding/json"
	"strconv"

	"github.com/jgr0sz/whoistrader/utils"
)

//Mapped JSON response from reverse.watch.
type ReverseWatchInfo struct {
	SteamID               string `json:"steam_id"`
	HasReversed           bool   `json:"has_reversed"`
	LastReversalTimestamp int64  `json:"last_reversal_timestamp,omitempty"`
}

//reverse.watch's endpoint.
type ReverseWatchEndpoint struct {
}

//Identifier for the map of retrieved data.
func (e *ReverseWatchEndpoint) Name() string {
	return "reverse.watch"
}

//Invokes and handles reverse.watch results.
func (e *ReverseWatchEndpoint) Fetch(steamID uint64) (any, error) {
	response, err := GetReverseWatchInfo(steamID)
	if err != nil {
		return "", err
	}
	return response, nil
}

//GETs and parses reverse.watch results.
func GetReverseWatchInfo(steamID uint64) (*ReverseWatchInfo, error) {
	url := "https://reverse.watch/api/v1/users/" + strconv.FormatUint(steamID, 10)
	body, err := utils.GetAPI(url, nil)
	if err != nil {
		return nil, err
	}

	var results ReverseWatchInfo
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}
	return &results, nil
}
