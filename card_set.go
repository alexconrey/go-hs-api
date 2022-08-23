package hs_api

import (
	"encoding/json"
	"fmt"
)

type CardSet struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (a *HearthstoneAPIClient) GetSets() ([]CardSet, error) {
	url := fmt.Sprintf("%s/metadata/sets?locale=%s", a.EndpointURL, a.Locale)
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return []CardSet{}, err
	}

	body, err := a.GetBodyFromRequest(req)
	if err != nil {
		return []CardSet{}, err
	}

	var cardSets []CardSet
	err = json.Unmarshal(body, &cardSets)
	if err != nil {
		return []CardSet{}, err
	}

	return cardSets, nil
}

func (a *HearthstoneAPIClient) GetSetByID(id int) (CardSet, error) {
	if len(a.CardSets) == 0 {
		return CardSet{}, fmt.Errorf("fatal: Card Sets not populated")
	}

	for _, cardSet := range a.CardSets {
		if cardSet.ID == id {
			return cardSet, nil
		}
	}
	return CardSet{}, nil
}
