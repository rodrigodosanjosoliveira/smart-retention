package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"smart-retention/internal/model"
)

type (
	CompraInput struct {
		ClienteID string            `json:"cliente_id"`
		Data      string            `json:"data"`
		Itens     []CompraItemInput `json:"itens"`
	}

	CompraItemInput struct {
		ItemID string  `json:"item_id"`
		Preco  float64 `json:"preco"`
	}

	CompraResponse struct {
		ClienteID   string         `json:"cliente_id"`
		NomeCliente string         `json:"nome_cliente"`
		Data        time.Time      `json:"data"`
		Itens       []ItemResponse `json:"itens"`
	}

	ItemResponse struct {
		Nome string `json:"nome"`
	}
)

func (h *Handler) CriarCompra(c *gin.Context) {
	var input CompraInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	data, err := time.Parse("2006-01-02", input.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "data inválida"})

		return
	}

	compra := model.Compra{
		ClienteID:  input.ClienteID,
		DataCompra: data,
	}

	tx := h.db.Begin()
	if err := tx.Create(&compra).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	for _, item := range input.Itens {
		compraItem := model.CompraItem{
			CompraID: compra.ID,
			ItemID:   item.ItemID,
			Preco:    item.Preco,
		}
		if err := tx.Create(&compraItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, compra)
}

func (h *Handler) ListarCompras(c *gin.Context) {
	var compras []model.Compra

	if err := h.db.Preload("Cliente").Preload("Itens.Item").Find(&compras).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	var response []CompraResponse
	for _, compra := range compras {
		var itens []ItemResponse
		for _, ci := range compra.Itens {
			itens = append(itens, ItemResponse{Nome: ci.Item.Nome})
		}
		response = append(response, CompraResponse{
			ClienteID:   compra.ClienteID,
			NomeCliente: compra.Cliente.Nome,
			Data:        compra.DataCompra,
			Itens:       itens,
		})
	}

	c.JSON(http.StatusOK, response)
}
