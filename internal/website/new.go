package website

import (
	"context"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
	"github.com/redis/go-redis/v9"
)

type Website struct {
	auth            AuthHandler
	serverInterface generated.ServerInterface
	htmlTemplater   templater.Templater
}

func NewWebsite(cfg WebsiteConfiguration) (*Website, error) {

	opt, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ac, err := NewAuthHandler(client, cfg)
	if err != nil {
		return nil, err
	}

	ops, err := storage.NewOperations(cfg.DatabaseUrl)
	if err != nil {
		return nil, err
	}
	err = ops.Migrate(context.Background())
	if err != nil {
		return nil, err
	}
	fsi := NewWebsiteWrapper(ops, cfg)

	tmpltr := templater.NewHtmlTemplater()
	return &Website{

		auth:            ac,
		htmlTemplater:   tmpltr,
		serverInterface: fsi,
	}, nil
}

type WebsiteWrapper struct {
	storage storage.Operations
	config  WebsiteConfiguration
}

func NewWebsiteWrapper(storage storage.Operations, cfg WebsiteConfiguration) generated.ServerInterface {
	return &WebsiteWrapper{
		storage: storage,
		config:  cfg,
	}
}
