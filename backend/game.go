package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Player struct {
	Username string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	/*
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
	*/
}

func makePlayer(email string) {

}
