package configs

type OAuth struct {
	ClientID            string   `config:"oauth_client_id"`
	ClientSecret        string   `config:"oauth_client_secret"`
	CookiesDomainName   string   `config:"oauth_cookies_domain_name"`
	FrontendRedirectURL string   `config:"oauth_frontend_redirect_url"`
	IssuerURL           string   `config:"oauth_issuer_url"`
	RedirectURL         string   `config:"oauth_redirect_url"`
	Scopes              []string `config:"oauth_scopes"`
}

func newOAuth() *OAuth {
	return &OAuth{
		CookiesDomainName:   "localhost",
		FrontendRedirectURL: "http://localhost:3000",
		RedirectURL:         "http://localhost:8080/v1/oauth/callback",
		Scopes:              []string{"profile", "email"},
	}
}
