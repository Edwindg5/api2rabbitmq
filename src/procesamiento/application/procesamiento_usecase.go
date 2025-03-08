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
	ch, err := p.RabbitConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declarar la cola (en caso de que no exista)
	_, err = ch.QueueDeclare(
		"pedido_enviado",
		true,  
		false, 
		false, 
		false, 
		nil,   
	)
	if err != nil {
		return err
	}

	// Convertir a JSON
	body, err := json.Marshal(pedido)
	if err != nil {
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

	return err
}
