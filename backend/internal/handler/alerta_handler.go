package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"smart-retention/internal/model"
)

type (
	Alerta struct {
		ClienteID      string   `json:"cliente_id"`
		NomeCliente    string   `json:"nome_cliente"`
		Motivo         string   `json:"motivo"`
		ItensFaltantes []string `json:"itens_faltantes,omitempty"`
	}

	AlertaResponse struct {
		ClienteID      string   `json:"cliente_id"`
		NomeCliente    string   `json:"nome_cliente"`
		Tipo           string   `json:"tipo"` // inatividade | item_faltando | dia_previsto
		Motivo         string   `json:"motivo"`
		ItensFaltantes []string `json:"itens_faltantes,omitempty"`
	}
)

func (h *Handler) GenerateAlerts() ([]Alerta, error) {
	var clientes []model.Cliente
	if err := h.db.Preload("DiasCompra").Preload("Itens").Find(&clientes).Error; err != nil {
		return nil, err
	}

	var alertas []Alerta
	hoje := time.Now().Weekday()

	for _, cliente := range clientes {
		deviaComprarHoje := false
		for _, dia := range cliente.DiasCompra {
			if int(hoje) == dia.DiaSemana {
				deviaComprarHoje = true
				break
			}
		}

		if !deviaComprarHoje {
			continue
		}

		// Verifica se o cliente comprou hoje
		var count int64
		h.db.Model(&model.Compra{}).
			Where("cliente_id = ? AND DATE(data_compra) = ?", cliente.ID, time.Now().Format("2006-01-02")).
			Count(&count)

		if count == 0 {
			alertas = append(alertas, Alerta{
				ClienteID:   cliente.ID,
				NomeCliente: cliente.Nome,
				Motivo:      "Cliente não comprou no dia esperado",
			})
			continue
		}

		// Verifica se comprou todos os itens
		var compra model.Compra
		err := h.db.Preload("Itens").
			Where("cliente_id = ? AND DATE(data_compra) = ?", cliente.ID, time.Now().Format("2006-01-02")).
			First(&compra).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		var itensCompradosHoje []model.CompraItem
		h.db.Where("compra_id = ?", compra.ID).Find(&itensCompradosHoje)

		mapaItensComprados := map[string]bool{}
		for _, item := range itensCompradosHoje {
			mapaItensComprados[item.ItemID] = true
		}

		var itensFaltando []string
		for _, item := range cliente.Itens {
			if !mapaItensComprados[item.ID] {
				itensFaltando = append(itensFaltando, item.Nome)
			}
		}

		if len(itensFaltando) > 0 {
			alertas = append(alertas, Alerta{
				ClienteID:      cliente.ID,
				NomeCliente:    cliente.Nome,
				Motivo:         "Itens não comprados na compra de hoje",
				ItensFaltantes: itensFaltando,
			})
		}
	}

	return alertas, nil
}

// GerarAlertasHoje Endpoint HTTP
func (h *Handler) GerarAlertasHoje(c *gin.Context) {
	alertas, err := h.GenerateAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alertas)
}

// DispararAlertaDiario Cron job
func (h *Handler) DispararAlertaDiario() {
	alertas, err := h.GenerateAlerts()
	if err != nil {
		log.Println("Erro ao gerar alertas:", err)
		return
	}

	for _, alerta := range alertas {
		log.Printf("⚠️ Alerta: %s - %s\n", alerta.NomeCliente, alerta.Motivo)
		if len(alerta.ItensFaltantes) > 0 {
			log.Printf("  Itens faltantes: %v\n", alerta.ItensFaltantes)
		}
	}
}

func (h *Handler) ListarAlertas(c *gin.Context) {
	var alertas []AlertaResponse
	diaAtual := int(time.Now().Weekday()) // 0 = Domingo

	// 1. Não comprou no dia previsto
	var clientes []model.Cliente
	h.db.Preload("DiasCompra").Find(&clientes)

	for _, cliente := range clientes {
		temDia := false
		for _, d := range cliente.DiasCompra {
			if d.DiaSemana == diaAtual {
				temDia = true
				break
			}
		}

		if temDia {
			var comprouHoje bool
			h.db.
				Model(&model.Compra{}).
				Where(
					`cliente_id = ? AND 
							data_compra >= (CURRENT_DATE AT TIME ZONE 'America/Sao_Paulo') AND 
							data_compra < ((CURRENT_DATE + INTERVAL '1 day') AT TIME ZONE 'America/Sao_Paulo')`,
					cliente.ID,
				).
				Select("count(*) > 0").
				Scan(&comprouHoje)

			if !comprouHoje {
				alertas = append(alertas, AlertaResponse{
					ClienteID:   cliente.ID,
					NomeCliente: cliente.Nome,
					Tipo:        "dia_previsto",
					Motivo:      "Hoje é um dia previsto e o cliente ainda não comprou.",
				})
			}
		}
	}

	// 2. Cliente inativo (última compra há mais de 7 dias)
	rows, err := h.db.Raw(`
	SELECT c.id, c.nome
	FROM clientes c
	LEFT JOIN compras co ON co.cliente_id = c.id
	GROUP BY c.id, c.nome
	HAVING MAX(co.data_compra) IS NULL OR MAX(co.data_compra) < CURRENT_DATE - INTERVAL '7 days'
`).Rows()

	if err != nil {
		log.Println("Erro ao executar SQL de inatividade:", err)

		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao buscar clientes inativos"})
	}
	defer rows.Close()

	for rows.Next() {
		var id, nome string
		if err := rows.Scan(&id, &nome); err != nil {
			log.Println("Erro ao escanear linha de inatividade:", err)
			continue
		}

		alertas = append(alertas, AlertaResponse{
			ClienteID:   id,
			NomeCliente: nome,
			Tipo:        "inatividade",
			Motivo:      "Cliente não compra há mais de 7 dias.",
		})
	}

	// 3. Cliente deixou de comprar itens que comprava
	h.db.Preload("Itens").Find(&clientes)
	for _, cliente := range clientes {
		var itensFaltantes []string

		for _, item := range cliente.Itens {
			var ultimaCompra time.Time
			h.db.Raw(`
				SELECT MAX(co.data_compra)
				FROM compra_items ci
				JOIN compras co ON co.id = ci.compra_id
				WHERE ci.item_id = ? AND co.cliente_id = ?
			`, item.ID, cliente.ID).Scan(&ultimaCompra)

			if ultimaCompra.Before(time.Now().AddDate(0, 0, -14)) {
				itensFaltantes = append(itensFaltantes, item.Nome)
			}
		}

		if len(itensFaltantes) > 0 {
			alertas = append(alertas, AlertaResponse{
				ClienteID:      cliente.ID,
				NomeCliente:    cliente.Nome,
				Tipo:           "item_faltando",
				Motivo:         "Cliente deixou de comprar itens recorrentes.",
				ItensFaltantes: itensFaltantes,
			})

			// Enviar alerta via WebSocket
			h.hub.BroadcastJSON(alertas)
		}
	}

	c.JSON(200, alertas)
}

func (h *Handler) BuildAllAlertas() ([]AlertaResponse, error) {
	var alertas []AlertaResponse
	diaAtual := int(time.Now().Weekday())

	// 1. Não comprou no dia previsto
	var clientes []model.Cliente
	if err := h.db.Preload("DiasCompra").Find(&clientes).Error; err != nil {
		return nil, err
	}

	for _, cliente := range clientes {
		temDia := false
		for _, d := range cliente.DiasCompra {
			if d.DiaSemana == diaAtual {
				temDia = true
				break
			}
		}
		if !temDia {
			continue
		}

		var comprouHoje bool
		h.db.
			Model(&model.Compra{}).
			Where("cliente_id = ? AND DATE(data_compra) = CURRENT_DATE", cliente.ID).
			Select("count(*) > 0").
			Scan(&comprouHoje)

		if !comprouHoje {
			alertas = append(alertas, AlertaResponse{
				ClienteID:   cliente.ID,
				NomeCliente: cliente.Nome,
				Tipo:        "dia_previsto",
				Motivo:      "Hoje é um dia previsto e o cliente ainda não comprou.",
			})
		}
	}

	// 2. Clientes inativos
	rows, err := h.db.Raw(`
		SELECT c.id, c.nome
		FROM clientes c
		LEFT JOIN compras co ON co.cliente_id = c.id
		GROUP BY c.id
		HAVING MAX(co.data_compra) IS NULL OR MAX(co.data_compra) < CURRENT_DATE - INTERVAL '7 days'
	`).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, nome string
		rows.Scan(&id, &nome)
		alertas = append(alertas, AlertaResponse{
			ClienteID:   id,
			NomeCliente: nome,
			Tipo:        "inatividade",
			Motivo:      "Cliente não compra há mais de 7 dias.",
		})
	}

	// 3. Itens deixados de comprar
	if err := h.db.Preload("Itens").Find(&clientes).Error; err != nil {
		return nil, err
	}

	for _, cliente := range clientes {
		var itensFaltantes []string
		for _, item := range cliente.Itens {
			var ultimaCompra time.Time
			h.db.Raw(`
				SELECT MAX(co.data_compra)
				FROM compra_items ci
				JOIN compras co ON co.id = ci.compra_id
				WHERE ci.item_id = ? AND co.cliente_id = ?
			`, item.ID, cliente.ID).Scan(&ultimaCompra)

			if ultimaCompra.Before(time.Now().AddDate(0, 0, -14)) {
				itensFaltantes = append(itensFaltantes, item.Nome)
			}
		}
		if len(itensFaltantes) > 0 {
			alertas = append(alertas, AlertaResponse{
				ClienteID:      cliente.ID,
				NomeCliente:    cliente.Nome,
				Tipo:           "item_faltando",
				Motivo:         "Cliente deixou de comprar itens recorrentes.",
				ItensFaltantes: itensFaltantes,
			})
		}
	}

	return alertas, nil
}
