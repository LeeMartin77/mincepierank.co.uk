package website

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/assets"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"

	"github.com/rs/zerolog/log"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/server.cfg.yaml ../../api/website.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/types.cfg.yaml ../../api/website.yaml

var (
	ErrShutdown error = fmt.Errorf("error shutting down server")
)

func handleStatic(w http.ResponseWriter, r *http.Request) {
	pth := filepath.Join("static", strings.TrimPrefix(r.URL.Path, "/static"))
	data, err := assets.StaticFS.ReadFile(pth)
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	mtype := mimetype.Lookup(filepath.Ext(pth))
	if mtype == nil {
		mtype = mimetype.Detect(data)
	}
	w.Header().Set("content-type", mtype.String())
	w.Write(data)
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	data, err := assets.StaticFS.ReadFile("static/favicon.ico")
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "image/x-icon")
	w.Write(data)
}

func (api Website) Run() error {
	r := gin.Default()

	r.GET("/static/*filepath", func(c *gin.Context) {
		handleStatic(c.Writer, c.Request)
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		handleFavicon(c.Writer, c.Request)
	})

	r.Use(api.auth.UserIdMiddleware())

	r.GET("/auth/signin", func(c *gin.Context) {
		api.auth.HandleSignin(c, c.Writer, c.Request)
	})

	r.GET("/auth/signout", func(c *gin.Context) {
		api.auth.HandleSignout(c)
	})

	r.GET("/auth/callback", func(c *gin.Context) {
		api.auth.HandleOAuth2Callback(c)
	})

	r.SetHTMLTemplate(api.htmlTemplater.GetTemplater())

	generated.RegisterHandlers(r, api.serverInterface)

	// Run the server
	port := os.Getenv("MPR_PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
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
