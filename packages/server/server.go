package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

func processMessage(client *Client, game *Game, message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("error unmarshaling message: %v", err)
		return
	}

	switch msg.Type {
	case StartGame:
		if client.Type == AdminClient {
			startGame(game)
		}

	case EndGame:
		if client.Type == AdminClient {
			endGame(game)
		}

	case AddScore:
		if client.Type == PlayerClient {
			game.mu.Lock()
			game.Score++
			game.mu.Unlock()
			updateAdminScore(game)
		}

	case LoseLife:
		if client.Type == PlayerClient {
			game.mu.Lock()
			game.Lives--
			game.mu.Unlock()
			updateAdminLives(game)
			if game.Lives <= 0 {
				game.IsActive = false
				notifyGameOver(game)
			}
		}

	}
}

func startGame(game *Game) {
	game.mu.Lock()
	game.IsActive = true
	game.Score = 0
	game.Lives = 3
	game.StartTime = time.Now()
	game.mu.Unlock()

	broadcastGameStart(game)

	go gameLoop(game)
}

func gameLoop(game *Game) {
	for {
		game.mu.RLock()
		isActive := game.IsActive
		game.mu.RUnlock()

		if !isActive {
			break
		}

		if time.Since(game.StartTime) > time.Second*60 {
			endGame(game)
			break
		}

		triggerRandomClientKeyChange(game)
		// TODO: Make configurable or based on the number of connected clients
		time.Sleep(time.Duration(rand.Float64())*time.Second*3 + time.Second)
	}
}

func triggerRandomClientKeyChange(game *Game) {
	game.mu.RLock()
	if len(game.PlayerClients) == 0 {
		game.mu.RUnlock()
		return
	}

	// TODO: See if there's a better way than iterating to strip this map
	clientIDs := make([]string, 0, len(game.PlayerClients))
	for id := range game.PlayerClients {
		clientIDs = append(clientIDs, id)
	}
	game.mu.RUnlock()

	randomIndex := rand.Intn(len(clientIDs))
	randomClientID := clientIDs[randomIndex]

	game.mu.Lock()
	client := game.PlayerClients[randomClientID]
	game.mu.Unlock()

	msg := Message{
		Type:    NewKey,
		Payload: json.RawMessage(`{}`),
	}

	msgBytes, _ := json.Marshal(msg)
	client.SendChan <- msgBytes
}
