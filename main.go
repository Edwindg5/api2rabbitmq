// api-notification//main.go
package main

import (
	"log"
	"net/http"
	"os"

	middleware "demo/src/core"
	"demo/src/procesamiento/application"
	"demo/src/procesamiento/infraestructure/controllers"
	"demo/src/procesamiento/infraestructure/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("❌ Error al cargar el archivo .env: %s", err)
	}

	// Obtener puerto de entorno
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("❌ No se encontró la variable de entorno PORT")
	}

	// Conectar a RabbitMQ
	rabbitMQURL := "amqp://admin:admin@52.7.35.94:5672/"
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("❌ Error al conectar con RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Inicializar caso de uso
	useCase := application.NewProcesadorPedidoUseCase(conn)

	// Configurar router y rutas
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API de Procesamiento de Pedidos"))
	})
	routes.RegisterProcesamientoRoutes(router, useCase)
	router.HandleFunc("/ws", controllers.WebSocketHandler)

	// Aplicar middleware de CORS
	server := &http.Server{
		Addr:    ":" + port,
		Handler: middleware.SetupCORS(router),
	}

	log.Println("🚀 Servidor escuchando en http://localhost:" + port)
	log.Fatal(server.ListenAndServe())
}
