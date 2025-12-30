package services

import (
	"sort"
	"sync"

	"github.com/raffaelramalhorosa/empacotador-api/models"
)

// logica de empacotar
type PackingService struct{}

// nova instancia de servico
func NewPackingService() *PackingService {
	return &PackingService{}
}

// PackOrders processa multiplas coisas
// pesquisando eu decidi tentar usar o Worker Pool para preparar o teste para 1000+ pedidos e consumir menos recursos da maquina
// "Por que você não usou uma rotina simples?"
// são 1 hora da manhã e eu fiquei curioso pra tentar, se der ruim volto pra versão sem channel e goRoutine
func (s *PackingService) PackOrders(orders []models.Order) []models.PackingResponse {

	numPedidos := len(orders)
	if numPedidos == 0 {
		return []models.PackingResponse{}
	}

	//aqui define o tanto de workers simultâneos
	numWorkers := 10
	if numPedidos < numWorkers {
		numWorkers = numPedidos
	}

	// canal que recebe os trabalhos
	jobCanal := make(chan jobData, numPedidos)

	//canal de resultados
	resultadoCanal := make(chan resultData, numPedidos)

	//isso aqui é pra sincronizar os workers
	var wg sync.WaitGroup

	//esse trecho todo é pra iniciar os workers
	for workerID := 1; workerID <= numWorkers; workerID++ {
		wg.Add(1)
		//cada um desses aqui é uma goutine que vai processar os pedidos
		go func(id int) {
			//esse trecho aqui marca como concluido se nao tiver mais nada a fazer
			defer wg.Done()

			// worker vai ficar em loop processando os pedidos do canal
			for job := range jobCanal {
				response := s.packSingleOrder(job.order)

				//vai enviar com o indice original
				resultadoCanal <- resultData{
					index:    job.index,
					response: response,
				}
			}
		}(workerID)
	}

	//vai enviar pro canal de trabalho, no pain no gain
	for i, order := range orders {
		jobCanal <- jobData{
			index: i,
			order: order,
		}
	}
	// essa coisa aqui é pra esperar os workers terminarem
	close(jobCanal)

	// goRoutine aqui pra fechar o canal quando todos os wkr terminarem
	go func() {
		wg.Wait()
		close(resultadoCanal)
	}()

	//coleta resultados do canal
	results := make([]resultData, 0, numPedidos)
	for result := range resultadoCanal {
		results = append(results, result)
	}

	//garantir que a ordem de tudo esta certa
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	responses := make([]models.PackingResponse, numPedidos)
	for i, result := range results {
		responses[i] = result.response
	}
	return responses
}

// isso aqui representa um work/pedido a ser processado
type jobData struct {
	index int
	order models.Order
}

// isso aqui representa m resultado processado
type resultData struct {
	index    int
	response models.PackingResponse
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

		// Se não cabe, cria nova caixa
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
