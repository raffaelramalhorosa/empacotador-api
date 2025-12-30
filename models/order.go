package models

type Order struct {
	ID       string    `json:"pedido_id"`
	Produtos []Product `json:"produtos"`
}
