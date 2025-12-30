package models

//reposta do empacotamento
type PackingResponse struct {
	PedidoID string `json:"pedido_id"`
	Caixas   []Box  `json:"caixas"`
}
