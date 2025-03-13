//api-notification//src/procesamiento/domain/entities/pedido.go
package entities

type Pedido struct {
	ID       int    `json:"id"`
	Cliente  string `json:"cliente"`
	Producto string `json:"producto"`
	Cantidad int    `json:"cantidad"`
	Estado   string `json:"estado"`
}
