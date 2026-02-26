package config

import (
	"os"
	"strings"
)

type Config struct {
	WeatherAPIKey     string
	ViaCEPBaseURL     string
	WeatherAPIBaseURL string
}

func Load() Config {
	_ = LoadDotEnv(".env")

	return Config{
		WeatherAPIKey:     strings.TrimSpace(os.Getenv("WEATHER_API_KEY")),
		ViaCEPBaseURL:     env("VIACEP_BASE_URL", "https://viacep.com.br"),
		WeatherAPIBaseURL: env("WEATHERAPI_BASE_URL", "https://api.weatherapi.com"),
	}
}

func env(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}
