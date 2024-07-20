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
	serverInterface generated.StrictServerInterface
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
	fsi := NewWebsiteWrapper(ops)

	return &Website{
		auth:            ac,
		serverInterface: fsi,
	}, nil
}

type WebsiteWrapper struct {
	storage       storage.Operations
	htmlTemplater templater.Templater
}

func NewWebsiteWrapper(storage storage.Operations) generated.StrictServerInterface {

	tmpltr := templater.NewHtmlTemplater()
	return &WebsiteWrapper{
		storage:       storage,
		htmlTemplater: tmpltr,
	}
}
