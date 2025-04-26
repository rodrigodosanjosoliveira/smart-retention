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
		CNPJ       string                   `json:"cnpj" binding:"required"`
		Nome       string                   `json:"nome" binding:"required"`
		Telefone   string                   `json:"telefone" binding:"required"`
		Email      string                   `json:"email"`
		Endereco   string                   `json:"endereco" binding:"required"`
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

func (h *Handler) AtualizarCliente(c *gin.Context) {
	clienteID := c.Param("id")

	var input model.Cliente
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	// Busca cliente original
	var cliente model.Cliente
	if err := h.db.Preload("Itens").Preload("DiasCompra").First(&cliente, "id = ?", clienteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Cliente não encontrado"})
		return
	}

	// Atualiza campos simples
	cliente.Nome = input.Nome
	cliente.CNPJ = input.CNPJ
	cliente.Telefone = input.Telefone
	cliente.Email = input.Email
	cliente.Endereco = input.Endereco

	// Atualiza campos simples no banco
	if err := h.db.Model(&cliente).Select("Nome", "CNPJ", "Telefone", "Email", "Endereco").Updates(cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao atualizar cliente"})
		return
	}

	// 🔁 Atualiza relação muitos-para-muitos: cliente_itens
	if err := h.db.Model(&cliente).Association("Itens").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao limpar itens do cliente"})
		return
	}
	if len(input.Itens) > 0 {
		if err := h.db.Model(&cliente).Association("Itens").Replace(input.Itens); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao atualizar itens do cliente"})
			return
		}
	}

	// 🔁 Atualiza dias de compra (relacionamento composto)
	if err := h.db.Where("cliente_id = ?", clienteID).Delete(&model.DiaCompraCliente{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao apagar dias de compra antigos"})
		return
	}
	if len(input.DiasCompra) > 0 {
		for i := range input.DiasCompra {
			input.DiasCompra[i].ClienteID = clienteID
		}
		if err := h.db.Create(&input.DiasCompra).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao salvar novos dias de compra"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Cliente atualizado com sucesso"})
}

func (h *Handler) DeletarCliente(c *gin.Context) {
	clienteID := c.Param("id")

	var cliente model.Cliente
	if err := h.db.Preload("Itens").Preload("DiasCompra").First(&cliente, "id = ?", clienteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Cliente não encontrado"})
		return
	}

	// Remove associação many2many explicitamente (cliente_itens)
	if err := h.db.Model(&cliente).Association("Itens").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao limpar itens associados"})
		return
	}

	// DiasCompra tem OnDelete:CASCADE, mas pode ser limpo explicitamente se quiser mais controle
	if err := h.db.Where("cliente_id = ?", cliente.ID).Delete(&model.DiaCompraCliente{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao remover dias de compra"})
		return
	}

	// Remove o cliente
	if err := h.db.Delete(&cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao excluir cliente"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) BuscarClientePeloID(c *gin.Context) {
	clienteID := c.Param("id")

	var cliente model.Cliente
	if err := h.db.Preload("Itens").Preload("DiasCompra").First(&cliente, "id = ?", clienteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Cliente não encontrado"})
		return
	}

	c.JSON(http.StatusOK, cliente)
}
