// api-notification//src/core/middleware.go
package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func SetupCORS(handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, 
		AllowCredentials: true,
	}).Handler(handler)
}


