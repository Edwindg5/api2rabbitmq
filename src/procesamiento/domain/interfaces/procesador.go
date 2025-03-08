//api-notification//src/procesamiento/domain/interfaces/procesador.go
package interfaces

import "demo/src/procesamiento/domain/entities"

type ProcesadorPedido interface {
	ProcesarPedido(pedido entities.Pedido) error
}
