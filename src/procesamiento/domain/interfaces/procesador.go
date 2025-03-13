// api-notification//src/procesamiento/domain/interfaces/procesador.go
package interfaces

import "demo/src/procesamiento/domain/entities"

type ProcesadorPedido interface {
	ProcesarPedido(pedido entities.Pedido) error
	EnviarPedidoEnviado(pedido entities.Pedido) error // Agrega esto si es parte de la lógica esperada
}
