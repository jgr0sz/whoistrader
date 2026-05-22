package endpoints

import (
	"encoding/json"
	"strconv"

	"github.com/jgr0sz/whoistrader/utils"
)

type ReverseWatchInfo struct {
	SteamID     string `json:"steam_id"`
	HasReversed bool   `json:"has_reversed"`
}

type ReverseWatchEndpoint struct {
}

func (e *ReverseWatchEndpoint) Name() string {
	return "reverse.watch"
}

func (e *ReverseWatchEndpoint) Fetch(steamID uint64) (any, error) {
	response, err := GetReverseWatchInfo(steamID)
	if err != nil {
		return "", err
	}
	return response, nil
}

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
