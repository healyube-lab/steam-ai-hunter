package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Game struct {
	ID       int     `json:"appid"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Discount float64 `json:"discount"`
}

var appIDs = []int{
	730, 570, 1172470, 271590, 1086940, 1245620, 292030, 1174180, 990080, 72850,
	550, 400, 252490, 346110, 227300, 218620, 252950, 440, 582010, 105600,
	230410, 220, 201510, 1515830, 413150, 412000, 256710, 1145460, 1046930, 238960,
}

func fetchPrice(appID int) (*Game, error) {
	resp, err := http.Get(fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%d&cc=us&l=en", appID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	appData, ok := data[fmt.Sprintf("%d", appID)].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no data for app %d", appID)
	}

	game := &Game{ID: appID}
	if d, ok := appData["data"].(map[string]interface{}); ok {
		game.Name, _ = d["name"].(string)
		if p, ok := d["price_overview"].(map[string]interface{}); ok {
			if final, ok := p["final"].(float64); ok {
				game.Price = final / 100.0
			}
			if discount, ok := p["discount_percent"].(float64); ok {
				game.Discount = discount
			}
		}
	}
	return game, nil
}

func main() {
	var games []*Game
	for _, id := range appIDs {
		game, err := fetchPrice(id)
		if err != nil {
			fmt.Printf("Error for %d: %v\n", id, err)
			time.Sleep(300 * time.Millisecond)
			continue
		}
		games = append(games, game)
		fmt.Printf("%s: $%.2f (discount: %.0f%%)\n", game.Name, game.Price, game.Discount)
		time.Sleep(300 * time.Millisecond)
	}

	out, _ := json.MarshalIndent(games, "", "  ")
	if err := os.WriteFile("../games_data.json", out, 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	}
}
