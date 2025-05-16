package controllers

import (
	"net/http"

	"github.com/H1vee/WeatherAPI/internal/services"
	"github.com/labstack/echo/v4"
)

type WeatherController struct {
	weatherService services.WeatherService
}

func NewWeatherController(weatherService services.WeatherService) *WeatherController {
	return &WeatherController{
		weatherService: weatherService,
	}
}

func (c *WeatherController) GetWeather(ctx echo.Context) error {
	city := ctx.QueryParam("city")
	if city == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "city parameter is required"})
	}
	weather, err := c.weatherService.GetCurrentWeather(city)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, weather)
}
