package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MessageType string

const (
	GameCreated        MessageType = "game_created"
	StartGame          MessageType = "start_game"
	EndGame            MessageType = "end_game"
	KeyPress           MessageType = "key_press"
	NewKey             MessageType = "new_key"
	AddScore           MessageType = "add_score"
	LoseLife           MessageType = "lose_life"
	UpdateScore        MessageType = "update_score"
	UpdateLives        MessageType = "update_lives"
	ClientJoined       MessageType = "client_joined"
	ClientDisconnected MessageType = "client_disconnected"

	AdminClient  ClientType = "admin"
	PlayerClient ClientType = "player"
)

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ClientType string

type Client struct {
	ID       string
	Type     ClientType
	Conn     *websocket.Conn
	SendChan chan []byte
}

type Game struct {
	ID            string
	AdminClient   *Client
	PlayerClients map[string]*Client
	Score         int
	Lives         int
	StartTime     time.Time
	IsActive      bool
	mu            sync.RWMutex
}

type GameManager struct {
	Games map[string]*Game
	mu    sync.RWMutex
}

var gameManager = &GameManager{
	Games: make(map[string]*Game),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: Fix for prod, allows all connections in development
	},
}

func main() {
	r := gin.Default()

	// Configure CORS to allow connections from Next.js
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/ws", handleWebSocket)
	r.GET("/games", func(c *gin.Context) {
		gameManager.mu.RLock()
		gameIDs := make([]string, 0, len(gameManager.Games))
		for id := range gameManager.Games {
			gameIDs = append(gameIDs, id)
		}
		gameManager.mu.RUnlock()

		c.JSON(http.StatusOK, gin.H{"games": gameIDs})
	})

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
