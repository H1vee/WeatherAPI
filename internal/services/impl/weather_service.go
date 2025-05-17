package impl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/H1vee/WeatherAPI/internal/services"
)

type weatherService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type weatherAPIResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`

	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  int     `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

func NewWeatherService(apiKey string) services.WeatherService {
	return &weatherService{
		apiKey:  apiKey,
		baseURL: "https://api.weatherapi.com/v1",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *weatherService) GetCurrentWeather(city string) (*services.WeatherData, error) {
	requestURL, err := url.Parse(fmt.Sprintf("%s/current.json", s.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	query := requestURL.Query()
	query.Set("key", s.apiKey)
	query.Set("q", city)
	requestURL.RawQuery = query.Encode()

	resp, err := s.httpClient.Get(requestURL.String())
	if err != nil {
		return nil, fmt.Errorf("weather API returned non-OK status: %d", resp.StatusCode)
	}

	var apiResp weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode weather API response: %w", err)
	}
	weatherData := &services.WeatherData{
		Temperature: apiResp.Current.TempC,
		Humidity:    apiResp.Current.Humidity,
		Description: apiResp.Current.Condition.Text,
	}
	return weatherData, nil
}
