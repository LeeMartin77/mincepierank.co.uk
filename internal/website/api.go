package website

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// convenience function for writing a JSON response
func writeApiJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error().Err(err).Msg("failed writing json response")
	}
}
