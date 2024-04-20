package website

import (
	"context"
	"io"
	"os"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

type Website struct {
	serverInterface generated.StrictServerInterface
}

func NewWebsite(out io.Writer) (*Website, error) {
	ops, err := storage.NewOperations(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	err = ops.Migrate(context.Background())
	if err != nil {
		return nil, err
	}
	fsi := NewWebsiteWrapper(ops)

	return &Website{
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
