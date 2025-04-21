package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"smart-retention/internal/model"
	"smart-retention/internal/ws"
	"time"
)

type (
	Handler struct {
		db  *gorm.DB
		hub *ws.Hub
	}

	ClienteInput struct {
		CNPJ       string                   `json:"cnpj"`
		Nome       string                   `json:"nome"`
		Telefone   string                   `json:"telefone"`
		Email      string                   `json:"email"`
		Endereco   string                   `json:"endereco"`
		Itens      []model.Item             `json:"itens"`
		DiasCompra []model.DiaCompraCliente `json:"dias_compra"`
	}
)

func NewHandler(db *gorm.DB, hub *ws.Hub) *Handler {
	return &Handler{
		db:  db,
		hub: hub,
	}
}

func (h *Handler) ListarClientes(c *gin.Context) {
	var clientes []model.Cliente

	if err := h.db.Preload("Itens").Preload("DiasCompra").Find(&clientes).Error; err != nil {
		c.JSON(500, gin.H{"error": "Erro ao listar clientes"})
		return
	}

	c.JSON(http.StatusOK, clientes)
}

func (h *Handler) CriarCliente(c *gin.Context) {
	var input ClienteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	cliente := model.Cliente{
		CNPJ:     input.CNPJ,
		Nome:     input.Nome,
		Telefone: input.Telefone,
		Email:    input.Email,
		Endereco: input.Endereco,
		Itens:    input.Itens,
	}

	// Cria cliente com itens
	if err := h.db.Create(&cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	// Remove dias de compra anteriores (por segurança, mesmo em criação inicial)
	if err := h.db.Where("cliente_id = ?", cliente.ID).Delete(&model.DiaCompraCliente{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	// Insere dias de compra com cliente_id correto
	for i := range input.DiasCompra {
		input.DiasCompra[i].ClienteID = cliente.ID
	}
	if len(input.DiasCompra) > 0 {
		if err := h.db.Create(&input.DiasCompra).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
			return
		}
	}

	// Retorna cliente com preload
	if err := h.db.Preload("Itens").Preload("DiasCompra").First(&cliente, "id = ?", cliente.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cliente)
}

func (h *Handler) HistoricoCliente(c *gin.Context) {
	clienteID := c.Param("id")

	var cliente model.Cliente
	if err := h.db.First(&cliente, "id = ?", clienteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Cliente não encontrado"})
		return
	}

	var compras []model.Compra
	h.db.Preload("Itens.Item").Where("cliente_id = ?", clienteID).Order("data_compra DESC").Find(&compras)

	type ItemCompra struct {
		Nome  string  `json:"nome"`
		Preco float64 `json:"preco"`
	}
	type CompraDTO struct {
		Data  time.Time    `json:"data"`
		Itens []ItemCompra `json:"itens"`
	}

	var historico []CompraDTO
	for _, compra := range compras {
		var itens []ItemCompra
		for _, ci := range compra.Itens {
			itens = append(itens, ItemCompra{
				Nome:  ci.Item.Nome,
				Preco: ci.Preco,
			})
		}
		historico = append(historico, CompraDTO{
			Data:  compra.DataCompra,
			Itens: itens,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"cliente": gin.H{
			"nome":     cliente.Nome,
			"cnpj":     cliente.CNPJ,
			"telefone": cliente.Telefone,
			"endereco": cliente.Endereco,
		},
		"historico": historico,
	})
}
