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
