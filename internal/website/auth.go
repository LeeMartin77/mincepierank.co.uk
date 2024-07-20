package website

import (
	"context"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
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
	HandleSignout(ctx *gin.Context)
	HandleOAuth2Callback(ctx *gin.Context)
	UserIdMiddleware() gin.HandlerFunc
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

func (ac *AuthConfig) HandleSignin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	state := uuid.NewString()
	resp := ac.Redis.Set(ctx, state, "true", time.Minute*15)
	if resp.Err() != nil {
		http.Error(w, resp.Err().Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, ac.Oauth2Config.AuthCodeURL(state), http.StatusFound)
}

// HandleSignout implements AuthHandler.
func (ac *AuthConfig) HandleSignout(ctx *gin.Context) {
	ctx.SetCookie(ac.CookieKey, "", -1, "/", ac.Config.MincepierankDomain, ac.Config.MincepierankSecureDomain, !ac.Config.MincepierankSecureDomain)
	ctx.Redirect(http.StatusFound, "/")
}

type jwtClaims struct {
	Email    string   `json:"email"`
	Verified bool     `json:"email_verified"`
	Groups   []string `json:"groups"`
}

func (ac *AuthConfig) HandleOAuth2Callback(ctx *gin.Context) {
	state := ctx.Request.URL.Query().Get("state")

	// Verify state.
	strd := ac.Redis.Get(ctx, state)

	if strd.Err() != nil || strd.Val() != "true" {
		if strd.Err() != nil {
			log.Error().Err(strd.Err()).Msg("Got an error trying to callback")
		} else {
			log.Warn().Any("resp", strd.Val()).Msg("non-true state value")
		}
		http.Error(ctx.Writer, "unrecognised state", http.StatusBadRequest)
		return
	}

	oauth2Token, err := ac.Oauth2Config.Exchange(ctx, ctx.Request.URL.Query().Get("code"))
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(ctx.Writer, "error reading token", http.StatusInternalServerError)
		return
	}

	// Parse and verify ID Token payload.
	idToken, err := ac.IdTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract custom claims.
	var claims jwtClaims
	if err := idToken.Claims(&claims); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	clms := jwt.MapClaims{
		"email": claims.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clms)

	tokenString, err := token.SignedString([]byte(ac.Config.JwtSigningKey))
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.SetCookie(ac.CookieKey, tokenString, 34560000, "/", ac.Config.MincepierankDomain, ac.Config.MincepierankSecureDomain, !ac.Config.MincepierankSecureDomain)
	ctx.Redirect(http.StatusFound, "/")
}

func (ac *AuthConfig) UserIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtStr, err := c.Cookie(ac.CookieKey)

		decodedToken, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(ac.Config.JwtSigningKey), nil
		})
		if err == nil && decodedToken != nil && decodedToken.Valid {
			// userid from token
			email := decodedToken.Claims.(jwt.MapClaims)["email"]
			if email != nil && email != "" {
				c.Keys = map[string]any{}
				c.Keys["userid"] = email
				c.Keys["signedin"] = true
				c.Next()
				return
			}
		}

		userid, err := c.Cookie("anon_user_id")

		if err != nil {
			userid = uuid.NewString()
			c.SetCookie("anon_user_id", userid, 34560000, "/", ac.Config.MincepierankDomain, ac.Config.MincepierankSecureDomain, !ac.Config.MincepierankSecureDomain)
		}
		c.Keys = map[string]any{}
		c.Keys["userid"] = userid
		c.Keys["signedin"] = false

		c.Next()
	}
}
