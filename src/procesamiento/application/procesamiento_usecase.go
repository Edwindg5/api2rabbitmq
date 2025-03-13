package application

import (
	"demo/src/procesamiento/domain/entities"
	"demo/src/procesamiento/domain/interfaces"
	"encoding/json"
	"log"


	"github.com/streadway/amqp"
)

type ProcesadorPedidoUseCase struct {
	Procesador interfaces.ProcesadorPedido
	RabbitConn *amqp.Connection
}

func NewProcesadorPedidoUseCase(rabbitConn *amqp.Connection) *ProcesadorPedidoUseCase {
	return &ProcesadorPedidoUseCase{
		RabbitConn: rabbitConn,
	}
}

func (p *ProcesadorPedidoUseCase) Procesar(pedido entities.Pedido) error {
	log.Printf("📦 Procesando pedido: %+v\n", pedido)

	// Enviar el pedido a la cola 'pedido_enviado'
	err := p.EnviarPedidoEnviado(pedido)
	if err != nil {
		log.Printf("❌ Error enviando pedido a 'pedido_enviado': %s\n", err)
		return err
	}

	log.Println("✅ Pedido enviado a la cola 'pedido_enviado'")
	return nil
}

func (p *ProcesadorPedidoUseCase) EnviarPedidoEnviado(pedido entities.Pedido) error {
	log.Printf("📤 Enviando pedido actualizado con estado: %s", pedido.Estado)

	ch, err := p.RabbitConn.Channel()
	if err != nil {
		log.Println("❌ Error al abrir canal RabbitMQ:", err)
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"pedido_enviado",
		true,  // Durable
		false, // AutoDelete
		false, // Exclusive
		false, // NoWait
		nil,
	)
	if err != nil {
		log.Println("❌ Error al declarar la cola 'pedido_enviado':", err)
		return err
	}

	body, err := json.Marshal(pedido)
	if err != nil {
		log.Println("❌ Error al serializar el pedido:", err)
		return err
	}

	err = ch.Publish(
		"",
		"pedido_enviado",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("❌ Error publicando mensaje en la cola 'pedido_enviado': %s", err)
	} else {
		log.Printf("✅ Mensaje enviado correctamente a la cola 'pedido_enviado': %s", string(body))
	}

	return err
}
