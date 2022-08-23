package hs_api

import (
	"encoding/json"
	"fmt"
)

type CardType struct {
	ID   int `json:"id"`
	Name string
	Slug string
}

func (a *HearthstoneAPIClient) GetCardTypes() ([]CardType, error) {
	url := fmt.Sprintf("%s/metadata/types?locale=%s", a.EndpointURL, a.Locale)
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return []CardType{}, err
	}

	body, err := a.GetBodyFromRequest(req)
	if err != nil {
		return []CardType{}, err
	}

	var cardTypes []CardType
	err = json.Unmarshal(body, &cardTypes)

	if err != nil {
		return []CardType{}, err
	}

	return cardTypes, nil
}

func (a *HearthstoneAPIClient) GetCardTypeByID(id int) (CardType, error) {
	if len(a.CardTypes) == 0 {
		return CardType{}, fmt.Errorf("fatal: Rarities not populated")
	}

	// Seek the correct rarity and return if it's found
	for _, cardType := range a.CardTypes {
		if cardType.ID == id {
			return cardType, nil
		}
	}
	return CardType{}, nil
}
