package handler

import (
	"github.com/gin-gonic/gin"
	"smart-retention/internal/model"
)

type DashboardResponse struct {
	TotalClientes      int                    `json:"total_clientes"`
	TotalCompras       int                    `json:"total_compras"`
	ComprasPorMes      []QuantidadePorPeriodo `json:"compras_por_mes"`
	ItensMaisComprados []QuantidadePorItem    `json:"itens_mais_comprados"`
	ClientesMaisAtivos []QuantidadePorCliente `json:"clientes_mais_ativos"`
}

type QuantidadePorPeriodo struct {
	Mes        string `json:"mes"`
	Quantidade int    `json:"quantidade"`
}

type QuantidadePorItem struct {
	Nome       string `json:"nome"`
	Quantidade int    `json:"quantidade"`
}

type QuantidadePorCliente struct {
	Nome       string `json:"nome"`
	Quantidade int    `json:"quantidade"`
}

func (h *Handler) ListarDashboard(c *gin.Context) {
	var totalClientes int64
	var totalCompras int64
	var comprasPorMes []QuantidadePorPeriodo
	var itensMaisComprados []QuantidadePorItem
	var clientesMaisAtivos []QuantidadePorCliente

	h.db.Model(&model.Cliente{}).Count(&totalClientes)
	h.db.Model(&model.Compra{}).Count(&totalCompras)

	// Compras por mês (YYYY-MM)
	h.db.
		Raw(`
			SELECT TO_CHAR(data_compra, 'YYYY-MM') as mes, COUNT(*) as quantidade
			FROM compras
			GROUP BY mes
			ORDER BY mes DESC
		`).Scan(&comprasPorMes)

	// Itens mais comprados
	h.db.
		Raw(`
			SELECT i.nome, COUNT(*) as quantidade
			FROM compra_items ci
			JOIN items i ON i.id = ci.item_id
			GROUP BY i.nome
			ORDER BY quantidade DESC
			LIMIT 5
		`).Scan(&itensMaisComprados)

	// Clientes mais ativos
	h.db.
		Raw(`
			SELECT c.nome, COUNT(*) as quantidade
			FROM compras co
			JOIN clientes c ON c.id = co.cliente_id
			GROUP BY c.nome
			ORDER BY quantidade DESC
			LIMIT 5
		`).Scan(&clientesMaisAtivos)

	res := DashboardResponse{
		TotalClientes:      int(totalClientes),
		TotalCompras:       int(totalCompras),
		ComprasPorMes:      comprasPorMes,
		ItensMaisComprados: itensMaisComprados,
		ClientesMaisAtivos: clientesMaisAtivos,
	}

	c.JSON(200, res)
}
