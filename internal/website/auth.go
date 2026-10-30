package website

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type AuthConfig struct {
	Oauth2Config    oauth2.Config
	IdTokenVerifier *oidc.IDTokenVerifier
	Redis           *redis.Client

	Config    WebsiteConfiguration
	CookieKey string
}

type AuthHandler interface {
	HandleSignin(ctx context.Context, w http.ResponseWriter, r *http.Request)
	HandleSignout(w http.ResponseWriter, r *http.Request)
	HandleOAuth2Callback(w http.ResponseWriter, r *http.Request)
	UserIdMiddleware() func(http.Handler) http.Handler
}

func NewAuthHandler(redis *redis.Client, config WebsiteConfiguration) (AuthHandler, error) {

	provider, err := oidc.NewProvider(context.Background(), config.AuthProviderUrl)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     config.AuthClientId,
		ClientSecret: config.AuthClientSecret,

		RedirectURL: config.AuthRedirectUrl,

		Endpoint: provider.Endpoint(),
		Scopes:   []string{oidc.ScopeOpenID, "profile", "email", "groups"},
	}

	// Create an ID token parser.
	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: config.AuthClientId})
	return &AuthConfig{
		Oauth2Config:    oauth2Config,
		IdTokenVerifier: idTokenVerifier,
		Redis:           redis,
		Config:          config,
		CookieKey:       "jwt",
	}, nil
}

type contextKey int

const (
	contextKeyUserID contextKey = iota
	contextKeySignedIn
)

// userIDFrom returns the id of the user attached to the request by
// UserIdMiddleware - their email when signed in, an anonymous id otherwise.
// It is empty if the middleware did not run.
func userIDFrom(r *http.Request) string {
	v, _ := r.Context().Value(contextKeyUserID).(string)
	return v
}

// signedInFrom reports whether the request carries a verified sign-in, as
// attached by UserIdMiddleware.
func signedInFrom(r *http.Request) bool {
	v, _ := r.Context().Value(contextKeySignedIn).(bool)
	return v
}

// setSessionCookie writes (or, when maxAge < 0, clears) one of the auth
// cookies. Secure domains get Secure cookies; HttpOnly mirrors the old gin
// behaviour, which set HttpOnly to the inverse of the secure flag.
func (ac *AuthConfig) setSessionCookie(w http.ResponseWriter, name, value string, maxAge int) {
	secure := ac.Config.MincepierankSecureDomain
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   ac.Config.MincepierankDomain,
		Secure:   secure,
		HttpOnly: !secure,
	})
}

// safeRedirectPath returns path if it is a safe same-site relative path to
// return to after signing in, or "" if it should be ignored.
func safeRedirectPath(path string) string {
	// Only allow single-slash absolute paths - blocks //host (protocol
	// relative URLs) and anything that could be an absolute URL. Backslashes
	// are rejected because browsers parse /\host like //host.
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") || strings.Contains(path, `\`) {
		return ""
	}
	if u, err := url.Parse(path); err != nil || u.Host != "" || u.IsAbs() {
		return ""
	}
	return path
}

// requireSignIn redirects the request to the sign-in route, asking to be
// returned to the requested path after signing in.
func requireSignIn(w http.ResponseWriter, r *http.Request) {
	target := "/auth/signin"
	if redirect := safeRedirectPath(r.URL.RequestURI()); redirect != "" {
		target = target + "?redirect=" + url.QueryEscape(redirect)
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func (ac *AuthConfig) HandleSignin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	state := uuid.NewString()

	// Persist where to return to after sign-in (if a safe path was given).
	redirect := safeRedirectPath(r.URL.Query().Get("redirect"))
	if redirect == "" {
		redirect = "/"
	}
	resp := ac.Redis.Set(ctx, state, redirect, time.Minute*15)
	if resp.Err() != nil {
		http.Error(w, resp.Err().Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, ac.Oauth2Config.AuthCodeURL(state), http.StatusFound)
}

// HandleSignout implements AuthHandler.
func (ac *AuthConfig) HandleSignout(w http.ResponseWriter, r *http.Request) {
	ac.setSessionCookie(w, ac.CookieKey, "", -1)
	http.Redirect(w, r, "/", http.StatusFound)
}

type jwtClaims struct {
	Email    string   `json:"email"`
	Verified bool     `json:"email_verified"`
	Groups   []string `json:"groups"`
}

func (ac *AuthConfig) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")

	// Verify state - the stored value is the path to return to after sign-in.
	strd := ac.Redis.Get(r.Context(), state)

	if strd.Err() != nil {
		log.Error().Err(strd.Err()).Msg("Got an error trying to callback")
		http.Error(w, "unrecognised state", http.StatusBadRequest)
		return
	}

	redirect := strd.Val()
	if safe := safeRedirectPath(redirect); safe != "" {
		redirect = safe
	} else {
		if redirect != "" {
			log.Warn().Str("state", state).Msg("discarding unsafe post-signin redirect")
		}
		redirect = "/"
	}

	oauth2Token, err := ac.Oauth2Config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "error reading token", http.StatusInternalServerError)
		return
	}

	// Parse and verify ID Token payload.
	idToken, err := ac.IdTokenVerifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract custom claims.
	var claims jwtClaims
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clms := jwt.MapClaims{
		"email": claims.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clms)

	tokenString, err := token.SignedString([]byte(ac.Config.JwtSigningKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ac.setSessionCookie(w, ac.CookieKey, tokenString, 34560000)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (ac *AuthConfig) UserIdMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			jwtCookie, err := r.Cookie(ac.CookieKey)
			if err == nil {
				decodedToken, err := jwt.Parse(jwtCookie.Value, func(token *jwt.Token) (interface{}, error) {
					return []byte(ac.Config.JwtSigningKey), nil
				})
				if err == nil && decodedToken != nil && decodedToken.Valid {
					// userid from token
					if email, ok := decodedToken.Claims.(jwt.MapClaims)["email"].(string); ok && email != "" {
						ctx = context.WithValue(ctx, contextKeyUserID, email)
						ctx = context.WithValue(ctx, contextKeySignedIn, true)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}

			userid := ""
			if anonCookie, err := r.Cookie("anon_user_id"); err == nil {
				userid = anonCookie.Value
			} else {
				userid = uuid.NewString()
				ac.setSessionCookie(w, "anon_user_id", userid, 34560000)
			}
			ctx = context.WithValue(ctx, contextKeyUserID, userid)
			ctx = context.WithValue(ctx, contextKeySignedIn, false)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
