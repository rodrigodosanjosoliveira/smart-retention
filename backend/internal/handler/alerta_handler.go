package handler

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"

	"smart-retention/internal/model"

	"github.com/gin-gonic/gin"
)

type (
	Alerta struct {
		ClienteID      string   `json:"cliente_id"`
		NomeCliente    string   `json:"nome_cliente"`
		Motivo         string   `json:"motivo"`
		ItensFaltantes []string `json:"itens_faltantes,omitempty"`
	}

	AlertaResponse struct {
		ClienteID       string          `json:"cliente_id"`
		NomeCliente     string          `json:"nome_cliente"`
		Tipo            string          `json:"tipo"` // inatividade | item_faltando | dia_previsto
		Motivo          string          `json:"motivo"`
		ItensFaltantes  []string        `json:"itens_faltantes,omitempty"`
		ItensDetalhados []ItemDetalhado `json:"itens_detalhados,omitempty"`
	}

	ItemDetalhado struct {
		Nome         string    `json:"nome"`
		UltimaCompra time.Time `json:"ultima_compra"`
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
	alertas, err := h.GerarTodosAlertas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	if len(alertas) > 0 {
		h.hub.BroadcastJSON(alertas)
	}

	c.JSON(http.StatusOK, alertas)
}

func (h *Handler) BuildAllAlertas() ([]AlertaResponse, error) {
	return h.GerarTodosAlertas()
}

func (h *Handler) GerarTodosAlertas() ([]AlertaResponse, error) {
	var alertas []AlertaResponse
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatalf("Erro ao carregar localização: %v", err)
	}

	diaAtual := int(time.Now().In(location).Weekday())

	var clientes []model.Cliente
	if err := h.db.Preload("DiasCompra").Preload("Itens").Find(&clientes).Error; err != nil {
		return nil, err
	}

	// 1. Não comprou no dia previsto
	for _, cliente := range clientes {
		deviaComprarHoje := false
		for _, d := range cliente.DiasCompra {
			if d.DiaSemana == diaAtual {
				deviaComprarHoje = true
				break
			}
		}

		if deviaComprarHoje {
			var comprouHoje bool
			h.db.
				Model(&model.Compra{}).
				Where(`
					cliente_id = ? AND 
					data_compra >= (CURRENT_DATE AT TIME ZONE 'America/Sao_Paulo') AND 
					data_compra < ((CURRENT_DATE + INTERVAL '1 day') AT TIME ZONE 'America/Sao_Paulo')
				`, cliente.ID).
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

	// 2. Clientes inativos
	rows, err := h.db.Raw(`
		SELECT c.id, c.nome
		FROM clientes c
		LEFT JOIN compras co ON co.cliente_id = c.id
		GROUP BY c.id, c.nome
		HAVING MAX(co.data_compra) IS NULL OR MAX(co.data_compra) < CURRENT_DATE - INTERVAL '7 days'
	`).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, nome string
		if err := rows.Scan(&id, &nome); err == nil {
			alertas = append(alertas, AlertaResponse{
				ClienteID:   id,
				NomeCliente: nome,
				Tipo:        "inatividade",
				Motivo:      "Cliente não compra há mais de 7 dias.",
			})
		}
	}

	// 3. Itens deixados de comprar
	for _, cliente := range clientes {
		var (
			itensFaltantes  []string
			itensDetalhados []ItemDetalhado
		)

		for _, item := range cliente.Itens {
			var tmp sql.NullTime
			row := h.db.Raw(`
				SELECT MAX(co.data_compra)
				FROM compra_items ci
				JOIN compras co ON co.id = ci.compra_id
				WHERE ci.item_id = ? AND co.cliente_id = ?
			`, item.ID, cliente.ID).Row()

			if err := row.Scan(&tmp); err != nil {
				log.Printf("Erro ao escanear última compra para item %s do cliente %s: %v", item.ID, cliente.ID, err)
				continue
			}

			if !tmp.Valid || tmp.Time.Before(time.Now().AddDate(0, 0, -14)) {
				itensFaltantes = append(itensFaltantes, item.Nome)

				itensDetalhados = append(itensDetalhados, ItemDetalhado{
					Nome:         item.Nome,
					UltimaCompra: tmp.Time,
				})
			}
		}

		if len(itensFaltantes) > 0 {
			alertas = append(alertas, AlertaResponse{
				ClienteID:       cliente.ID,
				NomeCliente:     cliente.Nome,
				Tipo:            "item_faltando",
				Motivo:          "Cliente deixou de comprar itens recorrentes.",
				ItensFaltantes:  itensFaltantes,
				ItensDetalhados: itensDetalhados,
			})
		}
	}

	return alertas, nil
}
