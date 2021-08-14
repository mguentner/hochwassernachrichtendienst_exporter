package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func HandleMetrics(w http.ResponseWriter, r *http.Request) {
	station := r.URL.Query().Get("station")
	if station == "" {
		http.Error(w, "Station ID not provided", http.StatusBadRequest)
		log.Info().Msg("Bad request: No Station ID provided")
		return
	}
	fetched, err := FetchLevels(station)
	if err != nil {
		log.Error().Msgf("Error during fetching station %s", station)
		if err == RequestError {
			http.Error(w, "GatewayError", http.StatusBadGateway)
			return
		} else if err == NotFoundError {
			http.Error(w, "NotFound", http.StatusNotFound)
			return
		} else {
			http.Error(w, "ServerError", http.StatusInternalServerError)
			return
		}
	}
	waterLevelInformation, err := ParseWaterLevel(fetched)
	if err != nil {
		log.Info().Msgf("ParserError: %s", err.Error())
		http.Error(w, "ParserError", http.StatusBadGateway)
		return
	}
	waterLevelInformation.StationID = station
	registry := prometheus.NewRegistry()
	err = waterLevelInformation.AddMetricsToRegistry(registry)
	if err != nil {
		log.Info().Msgf("ServerError: %s", err.Error())
		http.Error(w, "ServerError", http.StatusInternalServerError)
		return
	}
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
