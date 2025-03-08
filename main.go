package main

import (
	"demo/src/core/routes"
	"demo/src/procesamiento/application"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No se pudo cargar el archivo .env, verificando variables del sistema")
	}

	
	rabbitMQURL := "amqp://admin:admin@52.7.35.94:5672/"
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("❌ Error al conectar con RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Crear instancia del caso  con conexión a RabbitMQ
	useCase := application.NewProcesadorPedidoUseCase(conn)

	
	router := routes.SetupRouter(useCase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Println("🚀 Servidor corriendo en http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
