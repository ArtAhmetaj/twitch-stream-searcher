package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//TODO: check on generating these variables, take them through the webdriver if pure http requests will not cut it
var (
	baseUrl  = "https://gql.twitch.tv/gql"
	clientId = "kimne78kx3ncx6brgo4mv6wki5h1ko"
	hash     = "6ea6e6f66006485e41dbe3ebd69d5674c5b22896ce7b595d7fce6411a3790138"
)

type TwitchChannel struct {
	currentViews   int32
	channelLink    string
	similarityRate float32
}

func newTwitchChannel(value map[string]interface{}) TwitchChannel {
	return TwitchChannel{
		currentViews:   0,
		channelLink:    "",
		similarityRate: 0,
	}
}

type TwitchChannelSearchRequest struct {
	OperationName string `json:"operationName"`
	Variables     struct {
		Query     string  `json:"query"`
		Options   *string `json:"options"`
		RequestID *string `json:"requestID"`
	} `json:"variables"`
	Extensions struct {
		PersistedQuery struct {
			Version    int    `json:"version"`
			SHA256Hash string `json:"sha256Hash"`
		} `json:"persistedQuery"`
	} `json:"extensions"`
}

func newTwitchChannelSearchRequest(searchValue string) TwitchChannelSearchRequest {
	return TwitchChannelSearchRequest{
		OperationName: "SearchResultsPage_SearchResults",
		Variables: struct {
			Query     string  `json:"query"`
			Options   *string `json:"options"`
			RequestID *string `json:"requestID"`
		}{
			Query:     searchValue,
			Options:   nil,
			RequestID: nil,
		},
		Extensions: struct {
			PersistedQuery struct {
				Version    int    `json:"version"`
				SHA256Hash string `json:"sha256Hash"`
			} `json:"persistedQuery"`
		}{
			PersistedQuery: struct {
				Version    int    `json:"version"`
				SHA256Hash string `json:"sha256Hash"`
			}{
				Version:    1,
				SHA256Hash: hash,
			},
		},
	}
}

func getTwitchChannels(searchValue string) ([]TwitchChannel, error) {
	requestBody := newTwitchChannelSearchRequest(searchValue)
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(body))
	req.Header.Set("client-id", clientId)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	fmt.Println(req)
	client := &http.Client{}

	_, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	//TODO: add response struct and map to the smaller response that I need
	return nil, nil
}
