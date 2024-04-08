package website

import (
	"context"
	"fmt"
	"strings"

	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/server.cfg.yaml ../../website/website.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/types.cfg.yaml ../../website/website.yaml

var (
	ErrShutdown error = fmt.Errorf("error shutting down server")
)

func (wrpr *WebsiteWrapper) HomePage(_ context.Context, req generated.HomePageRequestObject) (generated.HomePageResponseObject, error) {

	str := "<h1>Hello, World!</h1>"
	reader := strings.NewReader(str)

	return generated.HomePage200TexthtmlResponse{
		Body:          reader,
		ContentLength: reader.Size(),
	}, nil
}
