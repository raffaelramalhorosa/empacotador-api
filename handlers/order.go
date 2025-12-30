package handlers

import (
	"net/http"

	"github.com/raffaelramalhorosa/empacotador-api/models"
	"github.com/raffaelramalhorosa/empacotador-api/services"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	packingService *services.PackingService
}

// NewOrderHandler cria uma nova instância do handler
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		packingService: services.NewPackingService(),
	}
}

// PackOrders godoc
// @Summary Empacotar pedidos
// @Description Recebe uma lista de pedidos e retorna o empacotamento otimizado
// @Tags Empacotamento
// @Accept json
// @Produce json
// @Param pedidos body []models.Order true "Lista de pedidos"
// @Success 200 {array} models.PackingResponse
// @Failure 400 {object} map[string]string
// @Router /api/pack [post]
func (h *OrderHandler) PackOrders(c *gin.Context) {
	var orders []models.Order

	// Faz bind do JSON recebido
	if err := c.ShouldBindJSON(&orders); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "JSON inválido: " + err.Error(),
		})
		return
	}

	// Valida se há pedidos
	if len(orders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Nenhum pedido foi enviado",
		})
		return
	}

	// Processa empacotamento
	responses := h.packingService.PackOrders(orders)

	c.JSON(http.StatusOK, responses)
}
