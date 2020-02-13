package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// room is a single chatroom
// We use channels to communicate concurrently between goroutines
// in a safe manner, as the room is to be shared amonst all clients.
type room struct {
	// forward is a channel that holds incoming messages
	// that are to be forward to other clients
	forward chan []byte
	// join is a channel for clients whishing to join teh room
	join chan *client
	// leave is a channel for clients wishing to leave
	leave chan *client
	// clients holds all current clients in this room
	clients map[*client]bool
}

// newRoom makes a new room.
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// forward the message to all clients
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// upgrade http connection into websocket
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	// create a client and setup a channel
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()

	// We kickoff client writes in a goroutine
	go client.write()

	// Reads happen in the main thread
	client.read()
}
