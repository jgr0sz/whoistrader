package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jgr0sz/whoistrader/utils"
	"github.com/tidwall/gjson"
)

//Response of json data
type CsFloatStats struct {
    MedianTradeTime     int `json:"median_trade_time"`
    TotalAvoidedTrades  int `json:"total_avoided_trades"`
    TotalFailedTrades   int `json:"total_failed_trades"`
    TotalTrades         int `json:"total_trades"`
    TotalVerifiedTrades int `json:"total_verified_trades"`
}

//Fulfilling Endpoint functions

type CSFloatEndpoint struct {
	APIKey string
}

func (e *CSFloatEndpoint) Name() string {
	return "csfloat"
}

func (e *CSFloatEndpoint) Fetch(steamID uint64) (any, error) {
	response, err := GetCsFloatStats(steamID, e.APIKey)
	if err != nil {
		return "", err
	}
	return response, nil
}

//GETs user CSFloat trade stats, which in turn will be funneled into the aggregator.
func GetCsFloatStats(steamID uint64, apiKey string) (*CsFloatStats, error){
	url := "https://csfloat.com/api/v1/users/" + strconv.FormatUint(steamID, 10) + "/stall"
	body, err := utils.GetAPI(url, map[string]string {
		"Authorization": apiKey,
	})
	if err != nil {
		return nil, err
	}

	userStats := gjson.GetBytes(body, "data.0.seller.statistics").Raw
	if userStats == "" {
		return nil, fmt.Errorf("User statistics not found (likely not an active seller?)")
	}
	
	var results CsFloatStats
	if err := json.Unmarshal([]byte(userStats), &results); err != nil {
		return nil, err
	}
	return &results, nil
}