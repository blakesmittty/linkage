package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Player struct {
	Username string
}

type GameState struct {
}

type GameRoom struct {
}

//rooms := make(GameRoom, 0)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	// for {
	// 	messageType, msg, err := ws.ReadMessage()
	// }
}

func makePlayer(email string) {

}
