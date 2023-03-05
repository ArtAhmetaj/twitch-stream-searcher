package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

//TODO: check on generating these variables, take them through the webdriver if pure http requests will not cut it
var (
	baseTwitchLink  = "https://www.twitch.tv"
	baseGraphqlLink = "https://gql.twitch.tv/gql"
	clientId        = "kimne78kx3ncx6brgo4mv6wki5h1ko"
	hash            = "6ea6e6f66006485e41dbe3ebd69d5674c5b22896ce7b595d7fce6411a3790138"
)

type TwitchChannel struct {
	displayName string
	followers   int
	channelLink string
}

func formatTwitchLinkByName(displayName string) string {
	return fmt.Sprintf("%s/%s", baseTwitchLink, displayName)
}

type TwitchChannelEdge struct {
	DisplayName string `json:"displayName"`
	Followers   struct {
		TotalCount int    `json:"totalCount"`
		Typename   string `json:"__typename"`
	} `json:"followers"`
	Description string `json:"description"`
	Channel     struct {
		ID       string `json:"id"`
		Typename string `json:"__typename"`
	} `json:"channel"`
	Stream *Stream `json:"stream"`
}

type Stream struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	ViewersCount int    `json:"viewersCount"`
	Typename     string `json:"__typename"`
}

type FreeformTag struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Typename string `json:"__typename"`
}

func parseEdges(response map[string]interface{}) ([]TwitchChannelEdge, error) {
	var items []TwitchChannelEdge
	dataNode := response["data"]
	searchFor, _ := dataNode.(map[string]interface{})["searchFor"]
	channels, _ := searchFor.(map[string]interface{})["channels"]
	edges, _ := channels.(map[string]interface{})["edges"]
	for _, e := range edges.([]interface{}) {
		var item TwitchChannelEdge
		//TODO: check on making it better, this is insanely ugly but is short code to convert a map to struct
		marshalledData, _ := json.Marshal(e.(map[string]interface{})["item"])
		err := json.Unmarshal(marshalledData, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func sortSliceFromSet(channels map[TwitchChannel]bool) []TwitchChannel {
	var channelsSlice []TwitchChannel
	for k := range channels {
		channelsSlice = append(channelsSlice, k)
	}
	sort.Slice(channelsSlice, func(i, j int) bool {
		return channelsSlice[i].followers > channelsSlice[j].followers
	})
	return channelsSlice
}

func GetTwitchChannels(edges []TwitchChannelEdge) []TwitchChannel {
	twitchChannels := map[TwitchChannel]bool{}
	for _, e := range edges {
		if e.Stream != nil {
			channel := TwitchChannel{
				displayName: e.DisplayName,
				followers:   e.Followers.TotalCount,
				channelLink: formatTwitchLinkByName(e.DisplayName),
			}
			twitchChannels[channel] = true
		}

	}
	return sortSliceFromSet(twitchChannels)
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

type TwitchClient struct {
	httpClient *http.Client
}

func NewTwitchClient() TwitchClient {
	return TwitchClient{
		httpClient: &http.Client{},
	}
}

func (tc TwitchClient) getTwitchChannels(searchValue string) ([]TwitchChannel, error) {
	requestBody := newTwitchChannelSearchRequest(searchValue)
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseGraphqlLink, bytes.NewBuffer(body))
	req.Header.Set("client-id", clientId)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	response, err := tc.httpClient.Do(req)

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}
	var parsedBody interface{}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return nil, err
	}

	parsedEdges, err := parseEdges(parsedBody.(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	return GetTwitchChannels(parsedEdges), nil
}
