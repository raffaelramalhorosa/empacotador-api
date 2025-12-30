package models

//Produto com suas dimensões e outras coisas importantes
type Product struct {
	ID          string  `json:"id"`
	Altura      float64 `json:"altura"`
	Largura     float64 `json:"largura"`
	Comprimento float64 `json:"comprimento"`
}

//trait pra calcular o quanto do produto ocupa de espaço na  caixinha
// p = this no javascript
func (p Product) Volume() float64 {
	return p.Altura * p.Largura * p.Comprimento
}
