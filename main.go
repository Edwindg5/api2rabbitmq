// api-notification//main.go
package main

import (
	middleware "demo/src/core"
	"demo/src/procesamiento/application"
	"demo/src/procesamiento/infraestructure/controllers"
	"demo/src/procesamiento/infraestructure/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("❌ Error al cargar el archivo .env: %s", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("❌ No se encontró la variable de entorno PORT")
	}

	rabbitMQURL := "amqp://admin:admin@52.7.35.94:5672/"
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("❌ Error al conectar con RabbitMQ: %s", err)
	}
	defer conn.Close()

	useCase := application.NewProcesadorPedidoUseCase(conn)

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API de Procesamiento de Pedidos"))
	})
	routes.RegisterProcesamientoRoutes(router, useCase)
	router.HandleFunc("/ws", controllers.WebSocketHandler)




	handler := middleware.SetupCORS(router)

	log.Println("Servidor escuchando en http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
