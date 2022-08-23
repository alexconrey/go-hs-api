package hs_api

import (
	"encoding/json"
	"fmt"
)

type Rarity struct {
	Slug        string
	ID          int
	CrafingCost []interface{} `json:"craftingCost"`
	DustValue   []interface{} `json:"DustValue"`
	Name        string
}

func (a *HearthstoneAPIClient) GetRarities() ([]Rarity, error) {
	url := fmt.Sprintf("%s/metadata/rarities?locale=%s", a.EndpointURL, a.Locale)
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return []Rarity{}, err
	}

	body, err := a.GetBodyFromRequest(req)
	if err != nil {
		return []Rarity{}, err
	}

	var rarities []Rarity
	err = json.Unmarshal(body, &rarities)

	if err != nil {
		return []Rarity{}, err
	}

	return rarities, nil

}

func (a *HearthstoneAPIClient) GetRarityById(id int) (Rarity, error) {
	if len(a.Rarities) == 0 {
		return Rarity{}, fmt.Errorf("fatal: Rarities not populated")
	}

	// Seek the correct rarity and return if it's found
	for _, rarity := range a.Rarities {
		if rarity.ID == id {
			return rarity, nil
		}
	}
	return Rarity{}, nil
}
