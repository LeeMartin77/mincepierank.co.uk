package website

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gabriel-vasile/mimetype"
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
	pth := "static" + strings.TrimPrefix(r.URL.Path, "/static")
	data, err := assets.StaticFS.ReadFile(pth)
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	ctype := mime.TypeByExtension(filepath.Ext(pth))
	if ctype == "" {
		ctype = mimetype.Detect(data).String()
	}
	w.Header().Set("content-type", ctype)
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

// statusRecorder remembers the status code written to the ResponseWriter so
// the logging middleware can report it.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs each request, replacing gin.Logger().
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", rec.status).
			Dur("took", time.Since(start)).
			Msg("request")
	})
}

// recoveryMiddleware turns panics into 500s, replacing gin.Recovery().
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error().Interface("panic", rec).Str("path", r.URL.Path).Msg("recovered from panic")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (api Website) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /static/{filepath...}", handleStatic)

	mux.HandleFunc("GET /favicon.ico", handleFavicon)

	mux.HandleFunc("GET /auth/signin", func(w http.ResponseWriter, r *http.Request) {
		api.auth.HandleSignin(r.Context(), w, r)
	})

	mux.HandleFunc("GET /auth/signout", api.auth.HandleSignout)

	mux.HandleFunc("GET /auth/callback", api.auth.HandleOAuth2Callback)

	// Routes from the OpenAPI spec get the userid middleware, as in the old
	// gin engine where r.Use applied it to everything registered after the
	// static and auth routes.
	apiMux := http.NewServeMux()
	generated.HandlerFromMux(api.serverInterface, apiMux)
	mux.Handle("/", api.auth.UserIdMiddleware()(apiMux))

	// Run the server
	port := os.Getenv("MPR_PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: loggingMiddleware(recoveryMiddleware(mux)),
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
