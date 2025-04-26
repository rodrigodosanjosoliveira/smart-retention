// Projeto: smart-retention
package main

import (
	"embed"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"io/fs"
	"log"
	"net/http"
	"smart-retention/internal/handler"
	"smart-retention/internal/infra/db"
	"smart-retention/internal/ws"
	"strings"
	"time"
)

//go:embed web/*
var embeddedFiles embed.FS

func main() {
	distFS, err := fs.Sub(embeddedFiles, "web")
	if err != nil {
		log.Fatal("Erro ao acessar arquivos embutidos:", err)
	}

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

	fileServer := http.FileServer(http.FS(distFS))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if !strings.Contains(path, ".") {
			r.URL.Path = "/index.html"
		}

		if r.URL.Path != "/index.html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}

		fileServer.ServeHTTP(w, r)
	})

	h := handler.NewHandler(dbConn, hub)

	r.GET("/clientes", h.ListarClientes)
	r.POST("/clientes", h.CriarCliente)
	r.POST("/compras", h.CriarCompra)
	r.GET("/alertas/hoje", h.GerarAlertasHoje)
	r.GET("/compras", h.ListarCompras)
	r.GET("/dashboard", h.ListarDashboard)
	r.GET("/alertas", h.ListarAlertas)
	r.GET("/ws/alertas", websocketHandler.HandleAlertasWS)
	r.GET("/clientes/:id/historico", h.HistoricoCliente)
	r.GET("/clientes/:id", h.BuscarClientePeloID)
	r.PUT("/clientes/:id", h.AtualizarCliente)
	r.DELETE("/clientes/:id", h.DeletarCliente)

	c := cron.New()

	c.AddFunc("0 8 * * *", func() {
		log.Println("📢 Rodando verificação de alertas automáticos...")
		h.DispararAlertaDiario()
	})

	c.Start()

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
