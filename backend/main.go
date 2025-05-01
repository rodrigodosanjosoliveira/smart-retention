// Projeto: smart-retention
package main

import (
	"log"
	"net/http"
	"smart-retention/internal/handler"
	"smart-retention/internal/infra/db"
	"smart-retention/internal/ws"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	dbConn := db.Connect()

	db.AutoMigrate(dbConn)

	r := gin.Default()

	hub := ws.NewHub()
	websocketHandler := &handler.WebSocketHandler{Hub: hub}
	alertaHandler := handler.NewHandler(dbConn, hub)

	go func() {
		for {
			time.Sleep(10 * time.Second)

			alertas, err := alertaHandler.BuildAllAlertas()
			if err != nil {
				log.Println("Erro ao gerar alertas:", err)
				continue
			}

			if len(alertas) > 0 {
				hub.BroadcastJSON(alertas)
			}
		}
	}()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h := handler.NewHandler(dbConn, hub)

	api := r.Group("/api")
	{
		api.GET("/clientes", h.ListarClientes)
		api.POST("/clientes", h.CriarCliente)
		api.POST("/compras", h.CriarCompra)
		api.GET("/alertas/hoje", h.GerarAlertasHoje)
		api.GET("/compras", h.ListarCompras)
		api.GET("/dashboard", h.ListarDashboard)
		api.GET("/alertas", h.ListarAlertas)
		api.GET("/ws/alertas", websocketHandler.HandleAlertasWS)
		api.GET("/clientes/:id/historico", h.HistoricoCliente)
		api.GET("/clientes/:id", h.BuscarClientePeloID)
		api.PUT("/clientes/:id", h.AtualizarCliente)
		api.DELETE("/clientes/:id", h.DeletarCliente)
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	c := cron.New()

	c.AddFunc("0 8 * * *", func() {
		log.Println("📢 Rodando verificação de alertas automáticos...")
		h.DispararAlertaDiario()
	})

	c.Start()

	err := r.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
