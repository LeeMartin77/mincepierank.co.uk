package website

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/server.cfg.yaml ../../website/website.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/types.cfg.yaml ../../website/website.yaml

var (
	ErrShutdown error = fmt.Errorf("error shutting down server")
)

func (wrpr *WebsiteWrapper) HomePage(c context.Context, req generated.HomePageRequestObject) (generated.HomePageResponseObject, error) {
	// urlExample := "postgres://username:password@localhost:5432/mincepierank"
	conn, err := pgx.Connect(c, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
		os.Exit(1)
	}
	defer conn.Close(c)

	var data string
	err = conn.QueryRow(c, "select test_data from test_table where id=$1", 1).Scan(&data)
	if err != nil {
		log.Error().Err(err).Msg("QueryRow failed")
		os.Exit(1)
	}

	str := fmt.Sprintf("<h1>Hello, World!</h1> <h3>%s</h3>", data)
	reader := strings.NewReader(str)

	return generated.HomePage200TexthtmlResponse{
		Body:          reader,
		ContentLength: reader.Size(),
	}, nil
}
