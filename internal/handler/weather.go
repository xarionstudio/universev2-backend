package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/pkg/response"
)

type WeatherHandler struct {
	client *http.Client
}

func NewWeatherHandler() *WeatherHandler {
	return &WeatherHandler{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *WeatherHandler) GetCurrentWeather(c fiber.Ctx) error {
	lat := c.Query("lat", "-0.5021")
	lon := c.Query("lon", "117.1536")

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,apparent_temperature,precipitation,weather_code,wind_speed_10m,is_day&timezone=auto",
		lat, lon,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create weather request: "+err.Error())
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return response.Error(c, fiber.StatusBadGateway, "Weather service unavailable: "+err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response.Error(c, resp.StatusCode, "Weather API error response")
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to parse weather data: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Success fetch weather from Open-Meteo", data)
}
