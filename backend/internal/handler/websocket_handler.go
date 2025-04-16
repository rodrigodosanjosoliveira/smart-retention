package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"smart-retention/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketHandler struct {
	Hub *ws.Hub
}

func (w *WebSocketHandler) HandleAlertasWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	w.Hub.AddClient(conn)

	for {
		// WebSocket clients don't need to send messages in this case.
		if _, _, err := conn.ReadMessage(); err != nil {
			w.Hub.RemoveClient(conn)
			break
		}
	}
}
