package services

import (
	"sort"

	"github.com/raffaelramalhorosa/empacotador-api/models"
)

// logica de empacotar
type PackingService struct{}

// nova instancia de servico
func NewPackingService() *PackingService {
	return &PackingService{}
}

// PackOrders processa multiplas coisas
func (s *PackingService) PackOrders(orders []models.Order) []models.PackingResponse {
	// boa prática de definir a capacidade com o tamanho exato da lista
	responses := make([]models.PackingResponse, 0, len(orders))

	// Loop
	for _, order := range orders {
		response := s.packSingleOrder(order)
		responses = append(responses, response)
	}

	return responses
}

// packSingleOrder empacota um único pedido
func (s *PackingService) packSingleOrder(order models.Order) models.PackingResponse {
	// Ordena produtos por volume
	produtos := make([]models.Product, len(order.Produtos))
	copy(produtos, order.Produtos)
	sort.Slice(produtos, func(i, j int) bool {
		return produtos[i].Volume() > produtos[j].Volume()
	})

	caixasUsadas := []models.Box{}

	// Tenta empacotar cada produto
	for _, produto := range produtos {
		empacotado := false

		// Tenta colocar em uma caixa existente
		for i := range caixasUsadas {
			if s.cabeProdutoNaCaixa(produto, caixasUsadas[i]) {
				caixasUsadas[i].Produtos = append(caixasUsadas[i].Produtos, produto)
				empacotado = true
				break
			}
		}

		// Se não coube, cria nova caixa
		if !empacotado {
			novaCaixa := s.selecionarMelhorCaixa(produto)
			if novaCaixa != nil {
				novaCaixa.Produtos = append(novaCaixa.Produtos, produto)
				caixasUsadas = append(caixasUsadas, *novaCaixa)
			}
		}
	}

	return models.PackingResponse{
		PedidoID: order.ID,
		Caixas:   caixasUsadas,
	}
}

// cabeProdutoNaCaixa verifica se o produto cabe nas dimensões da caixa
// Considera todas as rotações possíveis do produto
func (s *PackingService) cabeProdutoNaCaixa(produto models.Product, caixa models.Box) bool {
	// Calcula espaço ocupado pelos produtos já na caixa (simplificado)
	volumeOcupado := 0.0
	for _, p := range caixa.Produtos {
		volumeOcupado += p.Volume()
	}

	// Verifica se há volume disponível
	volumeDisponivel := caixa.Volume() - volumeOcupado
	if produto.Volume() > volumeDisponivel {
		return false
	}

	// Verifica se as dimensões cabem (considerando rotações)
	dimensoesProduto := []float64{produto.Altura, produto.Largura, produto.Comprimento}
	dimensoesCaixa := []float64{caixa.Altura, caixa.Largura, caixa.Comprimento}
	sort.Float64s(dimensoesProduto)
	sort.Float64s(dimensoesCaixa)

	// Compara dimensões ordenadas
	for i := 0; i < 3; i++ {
		if dimensoesProduto[i] > dimensoesCaixa[i] {
			return false
		}
	}

	return true
}

// selecionarMelhorCaixa escolhe a menor caixa que comporta o produto
func (s *PackingService) selecionarMelhorCaixa(produto models.Product) *models.Box {
	var melhorCaixa *models.Box

	for _, caixaTemplate := range models.TamanhosDisponiveis {
		// Cria uma cópia da caixa
		caixa := models.Box{
			ID:          caixaTemplate.ID,
			Altura:      caixaTemplate.Altura,
			Largura:     caixaTemplate.Largura,
			Comprimento: caixaTemplate.Comprimento,
			Produtos:    []models.Product{},
		}

		if s.cabeProdutoNaCaixa(produto, caixa) {
			if melhorCaixa == nil || caixa.Volume() < melhorCaixa.Volume() {
				melhorCaixa = &caixa
			}
		}
	}

	return melhorCaixa
}
