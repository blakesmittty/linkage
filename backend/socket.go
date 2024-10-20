package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { //allow all connections cuz who even hacks people
		return true
	},
}

var rooms = make(map[string]*GameRoom)
var mutex = &sync.Mutex{}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in socket handler")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	player := &Player{
		Username: r.URL.Query().Get("username"),
	}

	grid := &Grid{}
	initGrid(grid)

	roomID := r.URL.Query().Get("roomID")
	if roomID == "" {
		log.Println("Room ID not provided")
		return
	}

	fmt.Println("before first mutex lock")
	mutex.Lock()
	fmt.Println("after first mutext lock")
	room, exists := rooms[roomID]
	if !exists {
		fmt.Println("in if !exists")
		room = &GameRoom{
			ID:      roomID,
			Grid:    &Grid{},
			Players: make(map[*websocket.Conn]*Player),
		}
		fmt.Println("before init grid")
		initGrid(room.Grid)
		rooms[roomID] = room
	}
	fmt.Println("out of if !exists")
	room.Players[ws] = player
	mutex.Unlock()
	fmt.Println("mutex unlocked")

	fmt.Println("before broadcast")
	broadcastGrid(room)
	fmt.Println("after broadcast")

	fmt.Printf("before message loop")
	for {
		fmt.Println("in message loop")
		_, msg, err := ws.ReadMessage()
		fmt.Println("read message")
		fmt.Printf("msg: %v", msg)
		if err != nil {
			log.Println(err)
			removePlayerFromRoom(ws, roomID)
			return
		}

		log.Printf("ws message: %v\n", msg)

		var dropMsg BlockDropMessage
		err = json.Unmarshal(msg, &dropMsg)
		log.Printf("dropMsg: %v\n", dropMsg)
		if err != nil {
			log.Printf("error unmarshalling data for drop: %v\n", err)
			continue
		}

		if dropMsg.Action == "drop" {
			mutex.Lock()
			drop(room.Grid, dropMsg.Column, dropMsg.Block)
			mutex.Unlock()

			broadcastGrid(room)
		}

	}

}

func removePlayerFromRoom(conn *websocket.Conn, roomID string) {
	mutex.Lock()
	defer mutex.Unlock()

	room, exists := rooms[roomID]
	if exists {
		delete(room.Players, conn)
		conn.Close()

		// If no players remain in the room, delete the room
		if len(room.Players) == 0 {
			delete(rooms, roomID)
		}
	}
}

func broadcastGrid(room *GameRoom) {
	gridJSON, err := json.Marshal(room.Grid)
	if err != nil {
		log.Printf("Error marshaling grid to JSON: %v\n", err)
		return
	}

	for conn := range room.Players {
		err := conn.WriteMessage(websocket.TextMessage, gridJSON)
		if err != nil {
			log.Printf("Error sending grid message to player: %v\n", err)
		}
	}
}
