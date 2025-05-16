package services

type WeatherData struct {
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
	Description string  `json:"description"`
}

type WeatherService interface {
	GetCurrentWeather(city string) (*WeatherData, error)
}
