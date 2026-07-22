package handler

import (
	"github.com/gofiber/fiber/v3"
	"universev2-backend/pkg/response"
)

type WeatherHandler struct{}

func NewWeatherHandler() *WeatherHandler {
	return &WeatherHandler{}
}

func (h *WeatherHandler) GetCurrentWeather(c fiber.Ctx) error {
	// Mock proxy response from Open-Meteo
	data := fiber.Map{
		"current": fiber.Map{
			"temperature_2m": 31.5,
			"relative_humidity_2m": 65,
			"weather_code": 3,
			"wind_speed_10m": 12.4,
		},
		"hourly": fiber.Map{
			"time": []string{"2026-07-22T10:00", "2026-07-22T11:00"},
			"temperature_2m": []float64{31.5, 32.1},
			"precipitation_probability": []int{10, 20},
		},
	}
	return response.Success(c, fiber.StatusOK, "Success fetch weather", data)
}
