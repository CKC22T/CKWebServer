package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"packet"
	"room"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hashCount int = 0

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func eventHandle(c *Client, req *packet.RequestPacket) *packet.ResponsePacket {
	var res = packet.NewResponsePacket()
	switch req.Code {
	case packet.Login:
		res.Code = packet.Login
		log.Printf("Request Login")
		var nickname = req.Param["nickname"]
		if nickname == nil {
			//param 없음 오류
			res.Error = packet.Unknown
			break
		}
		res.Param["nickname"] = req.Param["nickname"]
		res.Param["hashTag"] = hashCount
		hashCount = hashCount + 1

	case packet.Logout:
		res.Code = packet.Logout
		log.Printf("Request Logout")
		c.hub.unregister <- c
	case packet.CreateRoom:
		res.Code = packet.CreateRoom
		log.Printf("Request CreateRoom")
		var r room.Room
		dediProc := room.DedicatedProcessOnBegin()
		r.Id = dediProc.Id
		r.Name = "temp"
		r.Addr = dediProc.Addr
		r.MaxUser = 4
		r.CurUser = 0
		res.Param["ip"] = r.Addr.Ip
		res.Param["port"] = r.Addr.Port
	case packet.LookUpRoom:
		res.Code = packet.LookUpRoom
		log.Printf("Request LookUpRoom")
		roomCode := int(req.Param["roomCode"].(float64))
		if roomCode == 0 {
			//param 없음 오류
			res.Error = packet.Unknown
		}
		var r = room.Rooms[roomCode]
		if r == nil {
			//방이 없음 오류
			res.Error = packet.Unknown
			r.Addr = room.Address{Ip: "", Port: 0}
		}
		res.Param["ip"] = r.Addr.Ip
		res.Param["port"] = r.Addr.Port
	case packet.Match:
		log.Printf("Request StartMatch")
		MatchHub.register <- c
	case packet.CancelMatch:
		res.Code = packet.CancelMatch
		log.Printf("Request CancelMatch")
		MatchHub.unregister <- c
	default:
		res.Error = packet.Unknown
	}

	return res
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		var res = packet.NewResponsePacket()
		if err != nil {
			res.Error = packet.Unknown
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		log.Println(string(message))

		var req packet.RequestPacket
		err = json.Unmarshal(message, &req)
		if err != nil {
			res.Error = packet.Unknown
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
		} else {
			res = eventHandle(c, &req)
		}

		if res.Code == packet.None {
			continue
		}

		buf, err := json.Marshal(res)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
		}

		c.send <- buf
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServerWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
