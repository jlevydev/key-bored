package main

import (
	"encoding/json"
	"fmt"
)

func broadcastGameStart(game *Game) {
	gameStateMsg := Message{
		Type: "game_state",
		Payload: json.RawMessage(`{"isActive":true,"score":` +
			fmt.Sprintf("%d", game.Score) + `,"lives":` +
			fmt.Sprintf("%d", game.Lives) + `}`),
	}

	msgBytes, _ := json.Marshal(gameStateMsg)

	game.AdminClient.SendChan <- msgBytes

	for _, client := range game.PlayerClients {
		client.SendChan <- msgBytes
	}
}

func notifyGameOver(game *Game) {
	gameOverMsg := Message{
		Type:    "game_over",
		Payload: json.RawMessage(`{"finalScore":` + fmt.Sprintf("%d", game.Score) + `}`),
	}

	msgBytes, _ := json.Marshal(gameOverMsg)

	// Send to admin
	game.AdminClient.SendChan <- msgBytes

	// Send to all players
	for _, client := range game.PlayerClients {
		client.SendChan <- msgBytes
	}
}

func disconnectClient(client *Client, game *Game) {
	client.Conn.Close()

	if client.Type == AdminClient {
		endGame(game)

		gameManager.mu.Lock()
		delete(gameManager.Games, game.ID)
		gameManager.mu.Unlock()
	} else {
		game.mu.Lock()
		delete(game.PlayerClients, client.ID)
		game.mu.Unlock()

		notifyAdminAboutPlayerDisconnect(game, client.ID)
	}
}

// TODO: Maybe organize this into a utils file since it's used elsewhere
func endGame(game *Game) {
	game.mu.Lock()
	game.IsActive = false
	game.mu.Unlock()

	notifyGameOver(game)
}
