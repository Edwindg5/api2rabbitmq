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
            // Filtra solo notificaciones válidas y en estado "pendiente"
            var pedidosPendientes []entities.Pedido
            for _, noti := range notificaciones {
                if noti.ID > 0 && noti.Cliente != "" && noti.Producto != "" && noti.Cantidad > 0 && noti.Estado == "pendiente" {
                    pedidosPendientes = append(pedidosPendientes, noti)
                }
            }
            json.NewEncoder(w).Encode(pedidosPendientes)
            return
        }

        var notificacion entities.Pedido
        if err := json.NewDecoder(r.Body).Decode(&notificacion); err != nil {
            http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
            return
        }

        // Asignar un ID correcto si es `0`
        if notificacion.ID == 0 {
            notificacion.ID = len(notificaciones) + 1
        }

        // Filtrar solo los pedidos con estado "pendiente"
        if notificacion.Estado != "pendiente" {
            log.Println("❌ Pedido rechazado: No está en estado 'pendiente'")
            http.Error(w, "Solo se aceptan pedidos pendientes", http.StatusBadRequest)
            return
        }

        log.Println("🔔 Nueva notificación recibida:", notificacion)

        // Almacenar en memoria
        notificaciones = append(notificaciones, notificacion)

        // Enviar a la cola 'pedido_enviado'
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
