package configs

import "time"

type HTTPServer struct {
	Addr                 string        `config:"http_server_addr"`
	CORSAllowedDomains   []string      `config:"http_server_cors_allowed_domains"`
	CORSAllowedMethods   []string      `config:"http_server_cors_allowed_methods"`
	CORSAllowedHeaders   []string      `config:"http_server_cors_allowed_headers"`
	CORSAllowCredentials bool          `config:"http_server_cors_allow_credentials"`
	CORSCacheMaxAge      time.Duration `config:"http_server_cors_cache_max_age"`
}

func newHTTPServer() *HTTPServer {
	return &HTTPServer{
		Addr: ":8080",
		CORSAllowedDomains: []string{
			"http://localhost:3000",
		},
		CORSAllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "HEAD", "OPTIONS"},
		CORSAllowedHeaders:   []string{"Authorization", "Origin", "Content-Length", "Content-Type"},
		CORSAllowCredentials: true,
		CORSCacheMaxAge:      12 * time.Hour,
	}
}
