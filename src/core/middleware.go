// api-notification//src/core/middleware.go
package middleware

import (
	"github.com/rs/cors"
	"net/http"
)

func SetupCORS(handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "ws://localhost:8081"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(handler)
	


}

