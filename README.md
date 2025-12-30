# API de Empacotamento de Pedidos

Sistema para empacotamento de produtos em caixas de papelão.

## Descrição

API REST desenvolvida em Go que recebe pedidos com produtos e suas dimensões, indicando quais produtos devem ser colocados em cada caixa. O sistema utiliza algoritmo First Fit Decreasing com suporte a rotação de produtos para minimizar o número de caixas utilizadas.

## Tecnologias Utilizadas

- Go 1.21+
- Gin Web Framework
- Swagger/OpenAPI para documentação
- Docker e Docker Compose
- GoRoutines e Channels para processamento concorrente

## Requisitos

### Desenvolvimento Local
- Go 1.21 ou superior
- Git

### Execução com Docker
- Docker 20.10+
- Docker Compose 1.29+

## Instalação

### Clonar Repositório

```bash
git clone https://github.com/raffaelramalhorosa/empacotador-api.git
cd empacotador-api
```

### Instalar Dependências

```bash
go mod download
go mod tidy
```

### Gerar Documentação Swagger

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

## Execução

### Modo Desenvolvimento

```bash
go run main.go
```

A API estará disponível em `http://localhost:8080`

### Modo Produção com Docker

```bash
# Build e iniciar
docker-compose up --build -d

# Ver logs
docker-compose logs -f

# Parar
docker-compose down
```

## Endpoints da API

### POST /api/pack

Processa pedidos e retorna empacotamento.

**Request:**
```json
[
  {
    "pedido_id": "pedido_001",
    "produtos": [
      {
        "id": "produto_1",
        "altura": 10,
        "largura": 15,
        "comprimento": 20
      }
    ]
  }
]
```

**Response:**
```json
[
  {
    "pedido_id": "pedido_001",
    "caixas": [
      {
        "caixa_id": 1,
        "altura": 30,
        "largura": 40,
        "comprimento": 80,
        "produtos": [
          {
            "id": "produto_1",
            "altura": 10,
            "largura": 15,
            "comprimento": 20
          }
        ]
      }
    ]
  }
]
```

### GET /health

Verifica status da aplicação.

**Response:**
```json
{
  "status": "ok"
}
```

### GET /swagger/index.html

Acessa documentação interativa da API.

## Caixas Disponíveis

O sistema trabalha com três tamanhos de caixas (altura x largura x comprimento em centímetros):

- Caixa 1: 30 x 40 x 80
- Caixa 2: 50 x 50 x 40
- Caixa 3: 50 x 80 x 60

## Algoritmo de Empacotamento

### Estratégia

O sistema utiliza o algoritmo First Fit Decreasing (FFD) com as seguintes características:

1. **Ordenação por Volume**: Produtos são ordenados por volume decrescente
2. **Tentativa de Reutilização**: Tenta encaixar produto em caixas já em uso
3. **Seleção Otimizada**: Escolhe a menor caixa que comporte o produto
4. **Rotação Automática**: Considera todas as orientações possíveis do produto
5. **Processamento Paralelo**: Utiliza Worker Pool com 10 GoRoutines para processar múltiplos pedidos simultaneamente

### Verificação de Encaixe

Para verificar se um produto cabe em uma caixa:

1. Verifica se o volume do produto cabe no espaço disponível
2. Ordena dimensões do produto e da caixa
3. Compara dimensões ordenadas (permite rotação)
4. Produto cabe se todas as dimensões forem menores ou iguais às da caixa

### Concorrência

O sistema implementa um Worker Pool que:

- Limita concorrência a 10 workers simultâneos
- Processa pedidos em paralelo usando GoRoutines
- Utiliza Channels para comunicação thread-safe
- Mantém ordem original dos pedidos na resposta
- Sincroniza workers usando WaitGroup

## Testes

### Teste Local

```bash
# Iniciar servidor
go run main.go

# Testar health check
curl http://localhost:8080/health

# Testar endpoint principal
curl -X POST http://localhost:8080/api/pack \
  -H "Content-Type: application/json" \
  -d '[{"pedido_id":"test","produtos":[{"id":"p1","altura":10,"largura":10,"comprimento":10}]}]'
```

### Teste com Docker

```bash
# Build e start
docker-compose up --build -d

# Aguardar inicialização (15 segundos)
sleep 15

# Testar
curl http://localhost:8080/health
```

### Swagger UI

Acesse `http://localhost:8080/swagger/index.html` para testes com interface gráfica.

## Estrutura do Projeto

```
empacotador-api/
├── models/              # Estruturas de dados
│   ├── box.go          # Definição de caixa e tamanhos disponíveis
│   ├── order.go        # Definição de pedido
│   ├── product.go      # Definição de produto
│   └── response.go     # Estrutura de resposta
├── services/            # Lógica de negócio
│   └── packing.go      # Algoritmo de empacotamento com Worker Pool
├── handlers/            # Controladores HTTP
│   └── order.go        # Handler do endpoint /api/pack
├── docs/               # Documentação Swagger (gerada automaticamente)
├── main.go             # Ponto de entrada da aplicação
├── Dockerfile          # Configuração Docker multi-stage
├── docker-compose.yml  # Orquestração de containers
├── go.mod              # Dependências Go
└── README.md           # Documentação
```

## Docker

### Build

```bash
docker-compose build
```

### Executar

```bash
# Foreground
docker-compose up

# Background
docker-compose up -d
```

### Ver Logs

```bash
docker-compose logs -f
```

### Parar

```bash
docker-compose down
```

### Verificar Status

```bash
docker ps
docker-compose logs api
```

## Características Técnicas

### Performance

- Processamento paralelo com Worker Pool
- Até 10 pedidos processados simultaneamente
- Multi-stage Docker build para imagens otimizadas

### Segurança

- Container executa com usuário não-root
- Health checks configurados
- Validação de entrada de dados

### Observabilidade

- Logs estruturados com Gin
- Health check endpoint
- Swagger para documentação e testes

## Autor

Rafael Ramalho Rosa

## Licença

Este projeto foi desenvolvido como parte de um teste técnico.