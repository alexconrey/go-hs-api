package hs_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Card struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	ClassID   int       `json:"classId"`
	CardClass CardClass `json:"class"`
	Mana      int       `json:"manaCost"`
	RarityID  int       `json:"rarityId"`
	Rarity    Rarity    `json:"rarity"`
	SetID     int       `json:"cardSetId"`
	Set       CardSet   `json:"cardSet"`
	TypeID    int       `json:"cardTypeId"`
	Type      CardType  `json:"type"`
}

type Cards struct {
	Cards     []Card
	PageCount int `json:"pageCount"`
	Page      int `json:"page"`
}

// Intention is for pagination as there are more than 1 page returned by a number of queries
func getRequestForPage(req *http.Request, page int, pageLimit int) *http.Request {
	q := req.URL.Query()
	// Remove any previous page values
	for _, key := range []string{"page", "pageLimit"} {
		for q.Get(key) != "" {
			q.Del(key)
		}
	}
	// Set the correct page value
	q.Add("page", strconv.FormatInt(int64(page), 10))
	q.Add("pageLimit", strconv.FormatInt(int64(pageLimit), 10))
	req.URL.RawQuery = q.Encode()
	return req
}

// This can be used for pagination with a loop and dynamic page set
func (a *HearthstoneAPIClient) GetCardsWithClassManaRaritySpec(class string, manaMin int, manaMax int, rarity string, page int, itemsPerPage int) ([]Card, error) {
	url := fmt.Sprintf("%s/cards", a.EndpointURL)
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return []Card{}, err
	}

	// Create query string for all possible mana values
	manaVals := []string{}
	i := manaMin
	for i <= manaMax {
		manaVals = append(manaVals, strconv.FormatInt(int64(i), 10))
		i += 1
	}
	manaStr := strings.Join(manaVals[:], ",")

	// Form the query
	q := req.URL.Query()
	q.Add("class", class)
	q.Add("locale", a.Locale)
	q.Add("manaCost", manaStr)
	q.Add("rarity", rarity)
	req.URL.RawQuery = q.Encode()

	pageReq := getRequestForPage(req, page, itemsPerPage)
	body, err := a.GetBodyFromRequest(pageReq)
	if err != nil {
		return []Card{}, err
	}

	var cards Cards
	err = json.Unmarshal(body, &cards)
	if err != nil {
		return []Card{}, err
	}

	if len(cards.Cards) == 0 {
		return []Card{}, fmt.Errorf("fatal: No cards found")
	}

	for idx, card := range cards.Cards {
		rarityByID, err := a.GetRarityById(card.RarityID)
		if err != nil {
			return []Card{}, err
		}
		cards.Cards[idx].Rarity = rarityByID

		cardSetByID, err := a.GetSetByID(card.SetID)
		if err != nil {
			return []Card{}, err
		}
		cards.Cards[idx].Set = cardSetByID

		cardClassByID, err := a.GetCardClassById(card.ClassID)
		if err != nil {
			return []Card{}, err
		}
		cards.Cards[idx].CardClass = cardClassByID

		cardTypeByID, err := a.GetCardTypeByID(card.TypeID)
		if err != nil {
			return []Card{}, err
		}
		cards.Cards[idx].Type = cardTypeByID
	}

	return cards.Cards, nil
}

func (a *HearthstoneAPIClient) GetCardsWithClassesManaRaritySpec(classes []string, manaMin int, manaMax int, rarity string) ([]Card, error) {
	var cards []Card
	var page = 1 // Only serve the first page per assignment requirements
	var itemsPerPage = 10
	for _, class := range classes {
		classCards, err := a.GetCardsWithClassManaRaritySpec(class, manaMin, manaMax, rarity, page, itemsPerPage)
		if err != nil {
			return []Card{}, err
		}
		cards = append(cards, classCards...)
	}
	return cards, nil
}
