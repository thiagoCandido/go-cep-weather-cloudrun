package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/example/go-cep-weather-cloudrun/internal/config"
	"github.com/example/go-cep-weather-cloudrun/internal/handler"
	"github.com/example/go-cep-weather-cloudrun/internal/services"
)

func main() {
	cfg := config.Load()

	client := &http.Client{Timeout: 8 * time.Second}

	via := services.NewViaCEP(client, cfg.ViaCEPBaseURL)
	wapi := services.NewWeatherAPI(client, cfg.WeatherAPIBaseURL, cfg.WeatherAPIKey)

	h := handler.New(handler.Deps{
		ViaCEP:  via,
		Weather: wapi,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// GET /weather?cep=01001000
	mux.Handle("/weather", h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr)

	// Cloud Run expects the process to listen on $PORT.
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
