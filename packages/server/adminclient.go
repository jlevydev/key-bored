package main

import (
	"encoding/json"
	"fmt"
)

// TODO: Maybe abstract these. Might be tough since some take different types as the raw json value
// But definitely worth looking at if we build out functionality
func notifyAdminAboutNewPlayer(game *Game, playerID string) {
	msg := Message{
		Type:    ClientJoined,
		Payload: json.RawMessage(`{"playerID":"` + playerID + `"}`),
	}

	msgBytes, _ := json.Marshal(msg)
	game.AdminClient.SendChan <- msgBytes
}

func updateAdminScore(game *Game) {
	msg := Message{
		Type:    UpdateScore,
		Payload: json.RawMessage(`{"score":` + fmt.Sprintf("%d", game.Score) + `}`),
	}

	msgBytes, _ := json.Marshal(msg)
	game.AdminClient.SendChan <- msgBytes
}

func updateAdminLives(game *Game) {
	msg := Message{
		Type:    UpdateLives,
		Payload: json.RawMessage(`{"lives":` + fmt.Sprintf("%d", game.Lives) + `}`),
	}

	msgBytes, _ := json.Marshal(msg)
	game.AdminClient.SendChan <- msgBytes
}

func notifyAdminAboutPlayerDisconnect(game *Game, playerID string) {
	msg := Message{
		Type:    ClientDisconnected,
		Payload: json.RawMessage(`{"playerID":"` + playerID + `"}`),
	}

	msgBytes, _ := json.Marshal(msg)
	game.AdminClient.SendChan <- msgBytes
}
