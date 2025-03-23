// api-notification//src/core/middleware.go
package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func SetupCORS(handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Permite cualquier origen
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Permite todos los encabezados
		AllowCredentials: true,
	}).Handler(handler)
}


