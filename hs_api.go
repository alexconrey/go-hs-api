package hs_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type BattleNetAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type HearthstoneAPIClient struct {
	AccessToken BattleNetAccessToken
	Rarities    []Rarity
	CardSets    []CardSet
	CardClasses []CardClass
	CardTypes   []CardType
	Locale      string
	httpClient  http.Client
	EndpointURL string
}

func (a *HearthstoneAPIClient) GetAccessToken(clientID string, clientSecret string) error {
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://us.battle.net/oauth/token", strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	// Basic Auth required
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Do the request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}

	// Wait for body to close before exiting function
	defer resp.Body.Close()

	// Get body as bytes object from response Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Try to unmarshal the body into a struct type
	var accessToken BattleNetAccessToken
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		fmt.Println(err.Error())
	}

	// If the AccessToken still doesn't exist, something bad has happened.
	if accessToken.AccessToken == "" {
		return fmt.Errorf("fatal: Access token not set despite a successful request: %s", string(body))
	}

	// If we're here, we've got a valid access token. No errors to report!
	a.AccessToken = accessToken

	return nil

}

func (a *HearthstoneAPIClient) DoRequest(req *http.Request) (*http.Response, error) {
	bearerHeader := fmt.Sprintf("Bearer %s", a.AccessToken.AccessToken)
	for req.Header.Get("Authorization") != "" {
		req.Header.Del("Authorization")
	}
	req.Header.Add("Authorization", bearerHeader)

	return a.httpClient.Do(req)
}

func (a *HearthstoneAPIClient) GetBodyFromRequest(req *http.Request) ([]byte, error) {
	resp, err := a.DoRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}

func (a *HearthstoneAPIClient) loadMetadata() error {
	rarities, err := a.GetRarities()
	if err != nil {
		return err
	}
	a.Rarities = rarities

	cardSets, err := a.GetSets()
	if err != nil {
		return err
	}
	a.CardSets = cardSets

	cardClasses, err := a.GetCardClasses()
	if err != nil {
		return err
	}
	a.CardClasses = cardClasses

	cardTypes, err := a.GetCardTypes()
	if err != nil {
		return err
	}
	a.CardTypes = cardTypes

	return nil

}

func NewRequest(verb string, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(verb, url, body)
}

func NewClient(clientID string, clientSecret string) (HearthstoneAPIClient, error) {
	hs_client := HearthstoneAPIClient{
		Locale:      "en_US",
		EndpointURL: "https://us.api.blizzard.com/hearthstone",
		httpClient:  http.Client{},
	}

	hs_client.GetAccessToken(clientID, clientSecret)

	// Populate metadata
	hs_client.loadMetadata()

	return hs_client, nil
}
