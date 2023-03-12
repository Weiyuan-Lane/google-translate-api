package http

import (
	nethttp "net/http"

	"github.com/gorilla/handlers"
)

func (h HttpServer) makeCORSWrappedHTTPHandler(handler nethttp.Handler) nethttp.Handler {
	corsHeaders := handlers.AllowedHeaders([]string{"x-requested-with", "origin"})
	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	return handlers.CORS(corsOrigins, corsHeaders, corsMethods)(handler)
}
