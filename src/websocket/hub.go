package websocket

import (
	"encoding/json"
	"log"
	"packet"
	"room"

	"github.com/gorilla/websocket"
)

var clientCount int = 0

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Matching() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			if len(h.clients) >= 2 {
				message := packet.NewResponsePacket()
				message.Code = packet.Match
				message.Error = packet.Success
				var r room.Room
				room.DedicatedProcessOnBegin(&r)
				r.Info.Name = "temp"
				message.Param["ip"] = r.Info.Addr.Ip
				message.Param["port"] = r.Info.Addr.Port
				buf, err := json.Marshal(message)
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
					}
				}
				for client := range h.clients {
					client.send <- buf
					delete(h.clients, client)
				}
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
				}
			}
		}
	}
}

var MatchHub *Hub

func (h *Hub) Run() {
	MatchHub = NewHub()
	go MatchHub.Matching()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			clientCount = clientCount + 1
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
				MatchHub.unregister <- client
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
