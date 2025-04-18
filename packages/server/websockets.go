package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}

	clientType := ClientType(c.Query("type"))
	// TODO: Handle if it does not exist. Currently creates an 'undefined' game room (spooky)
	gameID := c.Query("gameID")

	if clientType == AdminClient {
		registerAdminClient(conn, gameManager)
	} else {
		registerPlayerClient(conn, gameID, gameManager)
	}
}

func registerAdminClient(conn *websocket.Conn, gm *GameManager) {
	gameID := uuid.New().String()

	client := &Client{
		ID:       uuid.New().String(),
		Type:     AdminClient,
		Conn:     conn,
		SendChan: make(chan []byte, 256),
	}

	game := &Game{
		ID:            gameID,
		AdminClient:   client,
		PlayerClients: make(map[string]*Client),
		Score:         0,
		Lives:         3,
		IsActive:      false,
	}

	gm.mu.Lock()
	gm.Games[gameID] = game
	gm.mu.Unlock()

	msg := Message{
		Type:    GameCreated,
		Payload: json.RawMessage(`{"gameID":"` + gameID + `"}`),
	}

	msgBytes, _ := json.Marshal(msg)
	client.SendChan <- msgBytes

	go readPump(client, game)
	go writePump(client)
}

func registerPlayerClient(conn *websocket.Conn, gameID string, gm *GameManager) {
	gm.mu.RLock()
	game, exists := gm.Games[gameID]
	gm.mu.RUnlock()

	if !exists {
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Game not found"))
		conn.Close()
		return
	}

	client := &Client{
		ID:       uuid.New().String(),
		Type:     PlayerClient,
		Conn:     conn,
		SendChan: make(chan []byte, 256),
	}

	game.mu.Lock()
	game.PlayerClients[client.ID] = client
	game.mu.Unlock()

	notifyAdminAboutNewPlayer(game, client.ID)

	go readPump(client, game)
	go writePump(client)
}

func readPump(client *Client, game *Game) {
	defer func() {
		disconnectClient(client, game)
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		processMessage(client, game, message)
	}
}

func writePump(client *Client) {
	ticker := time.NewTicker(time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.SendChan:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.SendChan)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.SendChan)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
