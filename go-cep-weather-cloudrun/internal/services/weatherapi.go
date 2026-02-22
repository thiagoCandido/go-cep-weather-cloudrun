package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type WeatherAPI struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewWeatherAPI(client *http.Client, baseURL, apiKey string) *WeatherAPI {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	return &WeatherAPI{client: client, baseURL: baseURL, apiKey: strings.TrimSpace(apiKey)}
}

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (w *WeatherAPI) CurrentTempC(r *http.Request, city string, uf string) (float64, error) {
	if w.apiKey == "" {
		return 0, errors.New("WEATHER_API_KEY is required")
	}

	q := fmt.Sprintf("%s,%s,Brazil", city, uf)
	endpoint := w.baseURL + "/v1/current.json"

	u, _ := url.Parse(endpoint)
	qs := u.Query()
	qs.Set("key", w.apiKey)
	qs.Set("q", q)
	u.RawQuery = qs.Encode()

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather api status %d", resp.StatusCode)
	}

	var data weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Current.TempC, nil
}
