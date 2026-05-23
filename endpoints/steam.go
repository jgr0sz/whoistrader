package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/jgr0sz/whoistrader/utils"
	"github.com/tidwall/gjson"
)

const STEAM_API_DOMAIN = "https://api.steampowered.com/"

/* Steam exposes 3 useful endpoints for user information:
ISteamUser/GetPlayerSummaries, ISteamUser/GetPlayerBans, and IPlayerService/GetSteamLevel.
Since they are independent API calls, concurrently calling them would be much faster.


In the future, I plan to expand on this feature by adding an option to aggregate/check steam stats of
friends (since scammers tend to have dummy/malicious accounts added).
Will need to work on general caching/ratelimiting first.
*/

//Steam structs per API call
type SteamPlayerSummary struct {
	SteamID string `json:"steamid"`
	CommunityVisibilityState int `json:"communityvisibilitystate"`
	PersonaName string `json:"personaname"`
	TimeCreated int64 `json:"timecreated"`
}

type SteamPlayerBans struct {
	CommunityBanned bool `json:"CommunityBanned"`
	VACBanned bool `json:"VACBanned"`
	NumberOfVACBans int `json:"NumberOfVACBans"`
	DaysSinceLastBan int `json:"DaysSinceLastBan"`
	NumberOfGameBans int `json:"NumberOfGameBans"`
	EconomyBan string `json:"EconomyBan"`
}

type SteamLevel struct {
	PlayerLevel int `json:"player_level"`
}

//Overarching struct for the collected info
type SteamInfo struct {
	PlayerSummary *SteamPlayerSummary `json:"player_summary,omitempty"`
	PlayerBans *SteamPlayerBans `json:"player_bans,omitempty"`
	SteamLevel *SteamLevel	`json:"steam_level,omitempty"`
}

type SteamInfoEndpoints struct {
	APIKey string
}

func (e *SteamInfoEndpoints) Name() string {
	return "steamapi"
}

//GETs and parses player summary info.
func GetSteamPlayerSummary(steamID uint64, apiKey string) (*SteamPlayerSummary, error) {
	url :=  STEAM_API_DOMAIN + "ISteamUser/GetPlayerSummaries/v002/?key=" + apiKey + "&steamids=" + strconv.FormatUint(steamID, 10)
	body, err := utils.GetAPI(url, nil)
	if err != nil {
		return nil, err
	}

	playerSummary := gjson.GetBytes(body, "response.players.0").Raw
	if playerSummary == "" {
		return nil, fmt.Errorf("player summary info not found")
	}

	var results SteamPlayerSummary
	if err := json.Unmarshal([]byte(playerSummary), &results); err != nil {
		return nil, err
	}
	return &results, nil
}

//GETs and parses player ban info.
func GetSteamPlayerBans(steamID uint64, apiKey string) (*SteamPlayerBans, error) {
	url := STEAM_API_DOMAIN + "/ISteamUser/GetPlayerBans/v1/?key=" + apiKey + "&steamids=" + strconv.FormatUint(steamID, 10)
	body, err := utils.GetAPI(url, nil)
	if err != nil {
		return nil, err
	}

	playerBans := gjson.GetBytes(body, "players.0").Raw
	if playerBans == "" {
		return nil, fmt.Errorf("player bans info not found")
	}

	var results SteamPlayerBans
	if err := json.Unmarshal([]byte(playerBans), &results); err != nil {
		return nil, err
	}
	return &results, nil
}

//GETs and parses player level.
func GetSteamLevel(steamID uint64, apiKey string) (*SteamLevel, error){
	url := STEAM_API_DOMAIN + "IPlayerService/GetSteamLevel/v1/?key=" + apiKey + "&steamid=" + strconv.FormatUint(steamID, 10)
	body, err := utils.GetAPI(url, nil)
	if err != nil {
		return nil, err
	}

	playerLevel := gjson.GetBytes(body, "response").Raw
	if playerLevel == "" {
		return nil, fmt.Errorf("player level info not found")
	}

	var results SteamLevel
	if err := json.Unmarshal([]byte(playerLevel), &results); err != nil {
		return nil, err
	}
	return &results, nil
}

//Concurrently invokes and handles each of Steam web API calls, consolidating them into a SteamInfo struct.
//This is passed onto the main aggregator along with the information collected from the rest of the API calls in files in /endpoints/.
func (e *SteamInfoEndpoints) Fetch(steamID uint64) (any, error) {
	var (
		mutex sync.Mutex
		wg sync.WaitGroup
		info SteamInfo
		errors []string
	)

	wg.Go(func() {
		summary, err := GetSteamPlayerSummary(steamID, e.APIKey)
		mutex.Lock()
		defer mutex.Unlock()
		if err != nil {
			errors = append(errors, fmt.Sprintf("summary: %v", err))
			return
		}
		info.PlayerSummary = summary
	})

	wg.Go(func() {
        bans, err := GetSteamPlayerBans(steamID, e.APIKey)
        mutex.Lock()
        defer mutex.Unlock()
        if err != nil {
            errors = append(errors, fmt.Sprintf("bans: %v", err))
            return
        }
        info.PlayerBans = bans
    })

    wg.Go(func() {
        level, err := GetSteamLevel(steamID, e.APIKey)
        mutex.Lock()
        defer mutex.Unlock()
        if err != nil {
            errors = append(errors, fmt.Sprintf("level: %v", err))
            return
        }
        info.SteamLevel = level
    })
	wg.Wait()

	if info.PlayerSummary == nil && info.PlayerBans == nil && info.SteamLevel == nil {
		return nil, fmt.Errorf("all steam endpoints failed: %v", errors)
	}

	if len(errors) > 0 {
		log.Printf("Some sources failed for steamID %d: %v\n\n", steamID, errors)
	}
	return &info, nil
}
