package website

type WebsiteConfiguration struct {
	DatabaseUrl string `env:"DATABASE_URL, required"`
	RedisUrl    string `env:"REDIS_URL, required"`

	MincepierankDomain       string `env:"MINCEPIERANK_DOMAIN, required"`
	MincepierankSecureDomain bool   `env:"MINCEPIERANK_SECURE_DOMAIN, required"`

	AuthProviderUrl  string `env:"AUTH_PROVIDER_URL, required"`
	AuthClientId     string `env:"AUTH_CLIENT_ID, required"`
	AuthClientSecret string `env:"AUTH_CLIENT_SECRET, required"`
	AuthRedirectUrl  string `env:"AUTH_REDIRECT_URL, required"`
	JwtSigningKey    string `env:"JWT_SIGNING_KEY, required"`
}
