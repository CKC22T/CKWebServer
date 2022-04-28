package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"packet"
	"room"
	"strconv"
	"strings"
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

func logindb(nickname string, res *packet.ResponsePacket) {
	account := strings.Split(nickname, "#")
	req := packet.NewRequestPacket()

	uri := "signup"

	req.Param["nickname"] = account[0]
	if len(account) >= 2 {
		log.Printf("Login User : " + account[0] + " # " + account[1])
		req.Param["id"], _ = strconv.ParseInt(account[1], 10, 0)
		uri = "login"
	}
	request, _ := json.Marshal(req)
	reqBody := bytes.NewBufferString(string(request))
	resPost, err := http.Post("http://127.0.0.1:3010/"+uri, "text/plain", reqBody)
	if err != nil {
		res.Error = packet.Unknown
		panic(err)
	}
	defer resPost.Body.Close()

	json.NewDecoder(resPost.Body).Decode(res)
}

func eventLogin(res *packet.ResponsePacket, req *packet.RequestPacket) {
	log.Printf("Request Login")
	if !packet.ContainsParam(res, req, "nickname") {
		return
	}

	logindb(req.Param["nickname"].(string), res)
	if res.Error != packet.Success {
		return
	}

	res.Code = packet.Login
	res.Param["nickname"] = req.Param["nickname"]
	res.Param["hashTag"] = hashCount
	hashCount = hashCount + 1
}

func eventLogout(c *Client, res *packet.ResponsePacket, req *packet.RequestPacket) {
	log.Printf("Request Logout")
	res.Code = packet.Logout
	c.hub.unregister <- c
}

func eventCreateRoom(res *packet.ResponsePacket) {
	log.Printf("Request CreateRoom")
	res.Code = packet.CreateRoom
	var r room.Room
	room.DedicatedProcessOnBegin(&r)
	res.Param["ip"] = r.Info.Addr.Ip
	res.Param["port"] = r.Info.Addr.Port
}

func eventLookUpRoom(res *packet.ResponsePacket, req *packet.RequestPacket) {
	log.Printf("Request LookUpRoom")
	res.Code = packet.LookUpRoom
	if !packet.ContainsParam(res, req, "roomCode") {
		return
	}
	roomCode := int(req.Param["roomCode"].(float64))
	if roomCode == 0 {
		//param 없음 오류
		res.Error = packet.NotFoundParam
		return
	}
	var r = room.Rooms[roomCode]
	if r == nil {
		//방이 없음 오류
		res.Error = packet.NotFoundRoom
		return
	}
	res.Param["ip"] = r.Info.Addr.Ip
	res.Param["port"] = r.Info.Addr.Port
}

func eventMatch(c *Client) {
	log.Printf("Request StartMatch")
	MatchHub.register <- c
}

func eventCancelMatch(c *Client, res *packet.ResponsePacket) {
	log.Printf("Request CancelMatch")
	res.Code = packet.CancelMatch
	MatchHub.unregister <- c
}

func eventHandle(c *Client, req *packet.RequestPacket) *packet.ResponsePacket {
	var res = packet.NewResponsePacket()
	switch req.Code {
	case packet.Login:
		eventLogin(res, req)
	case packet.Logout:
		eventLogout(c, res, req)
	case packet.CreateRoom:
		eventCreateRoom(res)
	case packet.LookUpRoom:
		eventLookUpRoom(res, req)
	case packet.Match:
		eventMatch(c)
	case packet.CancelMatch:
		eventCancelMatch(c, res)
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
