// api-notification//src/procesamiento/infraestructure/controllers/procesamiento_controller.go
package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"demo/src/procesamiento/application"
	"demo/src/procesamiento/domain/entities"
)

func ProcesarPedido(useCase *application.ProcesadorPedidoUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pedido entities.Pedido
		json.NewDecoder(r.Body).Decode(&pedido)

		log.Println("📩 Pedido recibido en api-notifications:", pedido)
		err := useCase.Procesar(pedido)
		if err != nil {
			log.Println("❌ Error procesando pedido:", err)
			http.Error(w, "Error procesando pedido", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pedido)
	}
}


var notificaciones []entities.Pedido // Almacena las notificaciones en memoria

func NotificacionesHandler(useCase *application.ProcesadorPedidoUseCase) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet { // Si la petición es GET, devuelve las notificaciones almacenadas
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(notificaciones)
            return
        }

        var notificacion entities.Pedido
        if err := json.NewDecoder(r.Body).Decode(&notificacion); err != nil {
            http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
            return
        }

        log.Println("🔔 Nueva notificación recibida:", notificacion)

        // Guardar notificación en memoria para consultas futuras
        notificaciones = append(notificaciones, notificacion)

        // Enviar el mensaje a la cola 'pedido_enviado'
        err := useCase.EnviarPedidoEnviado(notificacion)
        if err != nil {
            log.Println("❌ Error enviando el mensaje a la cola 'pedido_enviado':", err)
            http.Error(w, "Error enviando el mensaje", http.StatusInternalServerError)
            return
        }

        log.Println("✅ Notificación almacenada y enviada correctamente")

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "Notificación enviada correctamente"})
    }
}

