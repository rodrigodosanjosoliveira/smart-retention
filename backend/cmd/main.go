// Projeto: smart-retention
package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"smart-retention/internal/handler"
	"smart-retention/internal/infra/db"
	"smart-retention/internal/ws"
	"time"
)

func main() {
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
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h := handler.NewHandler(dbConn, hub)

	r.GET("/clientes", h.ListarClientes)
	r.POST("/clientes", h.CriarCliente)
	r.POST("/compras", h.CriarCompra)
	r.GET("/alertas/hoje", h.GerarAlertasHoje)
	r.GET("/compras", h.ListarCompras)
	r.GET("/dashboard", h.ListarDashboard)
	r.GET("/alertas", h.ListarAlertas)
	r.GET("/ws/alertas", websocketHandler.HandleAlertasWS)

	c := cron.New()

	c.AddFunc("0 8 * * *", func() {
		log.Println("📢 Rodando verificação de alertas automáticos...")
		h.DispararAlertaDiario()
	})

	c.Start()

	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
