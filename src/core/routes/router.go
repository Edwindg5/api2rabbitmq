// api-notification//src/core/routes/router.go
package routes

import (

	"demo/src/procesamiento/application"
	"demo/src/procesamiento/infraestructure/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter(useCase *application.ProcesadorPedidoUseCase) *mux.Router {
	router := mux.NewRouter()
	

	routes.RegisterProcesamientoRoutes(router, useCase)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API de Procesamiento de Pedidos"))
	})

	return router
}
