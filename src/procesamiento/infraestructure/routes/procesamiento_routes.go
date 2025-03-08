package routes

import (
	"demo/src/procesamiento/application"
	"demo/src/procesamiento/infraestructure/controllers"

	"github.com/gorilla/mux"
)

func RegisterProcesamientoRoutes(router *mux.Router, useCase *application.ProcesadorPedidoUseCase) {
	router.HandleFunc("/procesar", controllers.ProcesarPedido(useCase)).Methods("POST")
}
