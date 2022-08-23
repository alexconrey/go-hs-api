package hs_api

import (
	"encoding/json"
	"fmt"
)

type CardClass struct {
	ID   int
	Name string
	Slug string
}

func (a *HearthstoneAPIClient) GetCardClasses() ([]CardClass, error) {
	url := fmt.Sprintf("%s/metadata/classes?locale=%s", a.EndpointURL, a.Locale)
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return []CardClass{}, err
	}

	body, err := a.GetBodyFromRequest(req)
	if err != nil {
		return []CardClass{}, err
	}

	var cardClasses []CardClass
	err = json.Unmarshal(body, &cardClasses)

	if err != nil {
		return []CardClass{}, err
	}

	return cardClasses, nil

}

func (a *HearthstoneAPIClient) GetCardClassById(id int) (CardClass, error) {
	if len(a.CardClasses) == 0 {
		return CardClass{}, fmt.Errorf("fatal: Rarities not populated")
	}

	// Seek the correct rarity and return if it's found
	for _, class := range a.CardClasses {
		if class.ID == id {
			return class, nil
		}
	}
	return CardClass{}, nil
}
