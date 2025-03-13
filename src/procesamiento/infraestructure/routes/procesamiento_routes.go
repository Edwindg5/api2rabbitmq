// api-notification//src/procesamiento/infraestructure/routes/procesamiento_routes.go
package routes

import (
	"demo/src/procesamiento/application"
	"demo/src/procesamiento/infraestructure/controllers"

	"github.com/gorilla/mux"
)

func RegisterProcesamientoRoutes(router *mux.Router, useCase *application.ProcesadorPedidoUseCase) {
	router.HandleFunc("/procesar", controllers.ProcesarPedido(useCase)).Methods("POST")
	router.HandleFunc("/notificaciones", controllers.NotificacionesHandler(useCase)).Methods("POST")
	router.HandleFunc("/ws", controllers.WebSocketHandler)
	router.HandleFunc("/notificaciones", controllers.NotificacionesHandler(useCase)).Methods("GET", "POST")


	

}
