package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	Username string
	Score    int
}

// GameState struct for storing important information about the game for single player or versus
type GameState struct {
	Grid [][]int
	Time time.Time
}

type GameRoom struct {
	ID      string
	State   GameState
	Players map[*websocket.Conn]string
}

//rooms := make(GameRoom, 0)

const GridHeight = 8
const GridWidth = 7

var grid [GridHeight][GridWidth]int

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { //allow all connections cuz who even hacks people
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	player := &Player{
		Username: r.URL.Query().Get("username"),
	}

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

	}

}

// func placeBlock() bool {

// }
