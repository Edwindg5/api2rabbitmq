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


var notificaciones []entities.Pedido // Lista en memoria

func NotificacionesHandler(useCase *application.ProcesadorPedidoUseCase) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(notificaciones) // Devuelve todos los pendientes
            return
        }

        var notificacion entities.Pedido
        if err := json.NewDecoder(r.Body).Decode(&notificacion); err != nil {
            http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
            return
        }

        if notificacion.Estado != "pendiente" {
            log.Println("❌ Pedido rechazado: Solo se aceptan pedidos 'pendiente'")
            http.Error(w, "Solo se aceptan pedidos pendientes", http.StatusBadRequest)
            return
        }

        log.Println("🔔 Notificación recibida y enviada a la cola pedido_enviado:", notificacion)

        // Enviar solo los pedidos "pendiente" a la cola 'pedido_enviado'
        err := useCase.EnviarPedidoEnviado(notificacion)
        if err != nil {
            log.Println("❌ Error enviando el mensaje a la cola:", err)
            http.Error(w, "Error enviando el mensaje", http.StatusInternalServerError)
            return
        }

        notificaciones = append(notificaciones, notificacion) // Guardar en memoria

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "Pedido enviado correctamente a la cola."})
    }
}
