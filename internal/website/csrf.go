package website

import (
	"context"
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	// csrfCookieName holds the random session id used to look up the
	// session's CSRF token in redis.
	csrfCookieName = "mpr_csrf_sid"
	// csrfHeader and csrfFormField are the two places a request can carry
	// the token back in: htmx sends the header (set on the body tag),
	// plain form posts send the hidden field.
	csrfHeader    = "X-CSRF-Token"
	csrfFormField = "csrf_token"

	csrfRedisPrefix = "csrf:"
	csrfTokenTTL    = 40 * 24 * time.Hour
)

// csrfStore keeps one CSRF token per session id. The redis-backed
// implementation matches the session storage used by auth; the interface
// keeps the middleware testable without a live redis.
type csrfStore interface {
	Get(ctx context.Context, sid string) (string, error)
	Set(ctx context.Context, sid string, token string) error
}

type redisCsrfStore struct {
	client *redis.Client
}

func (r redisCsrfStore) Get(ctx context.Context, sid string) (string, error) {
	token, err := r.client.Get(ctx, csrfRedisPrefix+sid).Result()
	if err == redis.Nil {
		return "", nil
	}
	return token, err
}

func (r redisCsrfStore) Set(ctx context.Context, sid string, token string) error {
	return r.client.Set(ctx, csrfRedisPrefix+sid, token, csrfTokenTTL).Err()
}

type csrfHandler struct {
	store  csrfStore
	config WebsiteConfiguration
}

func newCsrfHandler(client *redis.Client, config WebsiteConfiguration) *csrfHandler {
	return &csrfHandler{store: redisCsrfStore{client: client}, config: config}
}

// Middleware implements the synchroniser token pattern: every visitor gets a
// random session cookie, the matching token lives in redis and is handed to
// templates (via c.Keys["csrf_token"]) so it can ride along in forms and
// htmx headers. Every state-changing request must present it.
func (ch *csrfHandler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie(csrfCookieName)
		if err != nil || len(sid) == 0 || len(sid) > 64 {
			sid = uuid.NewString()
			c.SetCookie(csrfCookieName, sid, int(csrfTokenTTL.Seconds()), "/",
				ch.config.MincepierankDomain, ch.config.MincepierankSecureDomain, true)
		}

		token, err := ch.store.Get(c, sid)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if token == "" {
			token = uuid.NewString()
			if err := ch.store.Set(c, sid, token); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		if c.Keys == nil {
			c.Keys = map[any]any{}
		}
		c.Keys["csrf_token"] = token

		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			// Sliding expiry: keep the token alive while the session is in use.
			if err := ch.store.Set(c, sid, token); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			supplied := c.GetHeader(csrfHeader)
			if supplied == "" {
				supplied = c.Request.FormValue(csrfFormField)
			}
			if subtle.ConstantTimeCompare([]byte(supplied), []byte(token)) != 1 {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		c.Next()
	}
}
