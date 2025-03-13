// api-notification//src/procesamiento/infraestructure/controllers/websocket_controller.go
package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Gestor de conexiones WebSocket
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan interface{})

// Configuración del WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Maneja las conexiones WebSocket
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("❌ Error al conectar WebSocket:", err)
        return
    }
    defer conn.Close()

    clients[conn] = true
    log.Println("✅ Conexión WebSocket establecida")

    // Mantener conexión activa indefinidamente
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("❌ Error al leer mensaje del WebSocket:", err)
            delete(clients, conn)
            break
        }
        log.Printf("📩 Mensaje recibido: %s", string(msg))
    }
}

// Enviar mensaje a todos los clientes conectados
func BroadcastMessage(message interface{}) {
	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("❌ Error enviando mensaje por WebSocket:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
