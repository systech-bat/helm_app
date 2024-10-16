package iqair

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"weather-bot/pkg/model"

	log "github.com/sirupsen/logrus"
)

type IqairClient struct {
	apiURL   string
	apiToken string
}

func (a *IqairClient) GetAirQuality(country, state, city string) (*model.IqairResponse, error) {

	rawUrl, err := url.Parse(a.apiURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil, err
	}

	q := rawUrl.Query()
	q.Add("token", a.apiToken)

	rawUrl.RawQuery = q.Encode()
	requestURL := rawUrl.JoinPath("/city")
	v := url.Values{}
	v.Add("country", country)
	v.Add("state", state)
	v.Add("city", city)
	v.Add("key", a.apiToken)
	requestURL.RawQuery = v.Encode()

	log.Debug("Calling: ", requestURL.String())

	response, err := http.Get(requestURL.String())
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusBadRequest:
		return nil, fmt.Errorf("Bad Request")
	case http.StatusNotFound:
		return nil, fmt.Errorf("Page not found")
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Print the response body
	log.Debug("Successfully received payload: ", string(body))

	var iqairResponse *model.IqairResponse

	err = json.Unmarshal(body, &iqairResponse)
	if err != nil {
		return nil, err
	}
	return iqairResponse, nil
}

func NewIqairClient(url, token string) *IqairClient {
	return &IqairClient{
		apiURL:   url,
		apiToken: token,
	}
}

