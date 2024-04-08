package website

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"

	"github.com/rs/zerolog/log"
)

func (api Website) Run() error {
	// Create Gin API
	m := http.NewServeMux()

	foodHandler := generated.NewStrictHandler(api.serverInterface, nil)
	handler := generated.HandlerFromMux(foodHandler, m)

	// Run the server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8080"),
		Handler: handler,
	}

	log.Info().Msgf("starting on %s", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("critical error starting up")
		}
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	log.Info().Msg("gracefully shutting down")

	log.Info().Msg("server exiting")

	return nil
}
