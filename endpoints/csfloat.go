package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jgr0sz/whoistrader/utils"
	"github.com/tidwall/gjson"
)

//Mapped JSON response from CSFloat.
type CsFloatStats struct {
	MedianTradeTime     int `json:"median_trade_time"`
	TotalAvoidedTrades  int `json:"total_avoided_trades"`
	TotalFailedTrades   int `json:"total_failed_trades"`
	TotalTrades         int `json:"total_trades"`
	TotalVerifiedTrades int `json:"total_verified_trades"`
}

//Mapped Error response from CSFloat.
type CsFloatErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//CSFloat's endpoint.
type CSFloatEndpoint struct {
}

//Identifier for the map of retrieved data.
func (e *CSFloatEndpoint) Name() string {
	return "csfloat"
}

//Invokes and handles CSfloat stats.
func (e *CSFloatEndpoint) Fetch(ctx context.Context, steamID uint64) (any, error) {
	response, err := GetCsFloatStats(ctx, steamID)
	if err != nil {
		return nil, err
	}
	return response, nil
}

//GETs and parses user CSFloat trade stats.
func GetCsFloatStats(ctx context.Context, steamID uint64) (*CsFloatStats, error) {
	url := "https://csfloat.com/api/v1/users/" + strconv.FormatUint(steamID, 10) + "/stall"
	body, err := utils.GetAPI(ctx, url, nil)
	if err != nil {
		var csErr CsFloatErr
		if json.Unmarshal(body, &csErr) == nil && csErr.Message == "record not found" {
			return nil, fmt.Errorf("user has no CSFloat stall")
		}
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
