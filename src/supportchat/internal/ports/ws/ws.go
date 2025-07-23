package ws

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/MishaZhem/Gona/src/supportchat/internal/app"
	"github.com/MishaZhem/Gona/src/supportchat/internal/domain"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Handler struct {
	app         *app.App
	connections map[string]net.Conn
	mu          sync.Mutex
}

func NewHandler(app *app.App) *Handler {
	return &Handler{
		app:         app,
		connections: make(map[string]net.Conn),
	}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("user_id")
	role := r.URL.Query().Get("role")
	workerID := r.URL.Query().Get("worker_id")

	if userID == "" || (role != "user" && role != "worker") {
		wsutil.WriteServerMessage(conn, ws.OpText, []byte("missing or invalid parameters"))
		return
	}

	connID := userID
	if role == "worker" && workerID != "" {
		connID = workerID
	}
	h.mu.Lock()
	h.connections[connID] = conn
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.connections, connID)
		h.mu.Unlock()
	}()

	ctx := context.Background()

	for {
		msgBytes, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}

		text := strings.TrimSpace(string(msgBytes))
		if text == "" {
			continue
		}

		// to Redis
		var messages []*domain.Message
		switch role {
		case "user":
			messages, err = h.app.SendUserMessage(ctx, userID, text)
		case "worker":
			if workerID == "" {
				wsutil.WriteServerMessage(conn, ws.OpText, []byte("worker_id is required"))
				continue
			}
			messages, err = h.app.SendWorkerMessage(ctx, userID, workerID, text)
		}

		if err != nil {
			wsutil.WriteServerMessage(conn, ws.OpText, []byte("error: "+err.Error()))
			continue
		}

		// Other side
		var recipientID string
		if role == "user" {
			meta, err := h.app.GetChatMeta(ctx, userID)
			if err == nil && meta.WorkerID != "" {
				recipientID = meta.WorkerID
			}
		} else if role == "worker" {
			recipientID = userID
		}

		// Other
		if recipientID != "" {
			last := messages[len(messages)-1]
			msgJSON, _ := json.Marshal(last)

			// Sender
			wsutil.WriteServerMessage(conn, ws.OpText, msgJSON)

			h.mu.Lock()
			targetConn, ok := h.connections[recipientID]
			h.mu.Unlock()
			if ok {
				_ = wsutil.WriteServerMessage(targetConn, ws.OpText, msgJSON)
			}
		} else {
			msgWithBotJSON, _ := json.Marshal(messages[max(len(messages)-2, 0):])
			wsutil.WriteServerMessage(conn, ws.OpText, msgWithBotJSON)
		}
	}
}
