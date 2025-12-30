package models

//caixa de papelão e dimensoes
//antes era json:produtos.altura mas descobri q o GO não entende aninhamentos assim
type Box struct {
	ID          int       `json:"id"`
	Altura      float64   `json:"altura"`
	Largura     float64   `json:"largura"`
	Comprimento float64   `json:"comprimento"`
	Produtos    []Product `json:"produtos"`
}

func (b Box) Volume() float64 {
	return b.Altura * b.Largura * b.Comprimento
}

//os tamanhos disponúveis
//descobri q da para chamar models.TamanhosDisp globalmente
var TamanhosDisponiveis = []Box{
	{ID: 1, Altura: 30, Largura: 40, Comprimento: 80},
	{ID: 2, Altura: 50, Largura: 50, Comprimento: 40},
	{ID: 3, Altura: 50, Largura: 80, Comprimento: 60},
}
