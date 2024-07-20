package main

import (
	"context"
	"errors"

	"github.com/leemartin77/mincepierank.co.uk/internal/website"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

func main() {

	var config website.WebsiteConfiguration
	if err := envconfig.Process(context.Background(), &config); err != nil {
		log.Fatal().Err(err).Msg("An error was encountered parsing environment variables")
	}

	site, err := website.NewWebsite(config)

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
