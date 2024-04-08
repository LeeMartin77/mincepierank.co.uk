package main

import (
	"errors"
	"os"

	"github.com/leemartin77/mincepierank.co.uk/internal/website"
	"github.com/rs/zerolog/log"
)

func main() {

	site, err := website.NewWebsite(os.Stdout)

	if err != nil {
		log.Fatal().Err(err).Msg("api failed to initialise")
	}

	if err := site.Run(); err != nil {
		if errors.Is(err, website.ErrShutdown) {
			log.Error().Err(err).Msg("api failed to shutdown")
		} else {
			log.Fatal().Err(err).Msg("api failed to start")
		}
	}
}
