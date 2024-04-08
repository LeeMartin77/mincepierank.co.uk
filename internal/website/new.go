package website

import (
	"io"

	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

type Website struct {
	serverInterface generated.StrictServerInterface
}

func NewWebsite(out io.Writer) (*Website, error) {

	fsi := NewWebsiteWrapper()

	return &Website{
		serverInterface: fsi,
	}, nil
}

type WebsiteWrapper struct{}

func NewWebsiteWrapper() generated.StrictServerInterface {
	return &WebsiteWrapper{}
}
