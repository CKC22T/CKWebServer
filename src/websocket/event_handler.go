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

	"github.com/gorilla/websocket"
)

var dbUrl string = "http://127.0.0.1:3010"

func requestDB(url string, res *packet.ResponsePacket, req *packet.RequestPacket) {
	request, _ := json.Marshal(req)
	reqBody := bytes.NewBufferString(string(request))
	resPost, err := http.Post(dbUrl+url, "text/plain", reqBody)
	if err != nil {
		res.Error = packet.Unknown
		log.Print(err.Error())
		return
	}
	defer resPost.Body.Close()

	json.NewDecoder(resPost.Body).Decode(res)
}

func logindb(nickname string, id int, res *packet.ResponsePacket) {
	req := packet.NewRequestPacket()
	req.Param["nickname"] = nickname
	req.Param["id"] = id
	requestDB("/login", res, req)
}

func signupdb(nickname string, res *packet.ResponsePacket) {
	req := packet.NewRequestPacket()
	req.Param["nickname"] = nickname
	requestDB("/signup", res, req)
}

func eventLogin(res *packet.ResponsePacket, req *packet.RequestPacket) {
	log.Printf("Request Login")
	if !packet.ContainsParamReq(res, req, "nickname") {
		return
	}

	account := strings.Split(req.Param["nickname"].(string), "#")

	nickname := account[0]
	id := 0

	//req.Param["nickname"] = account[0]
	if len(account) >= 2 {
		id, _ := strconv.ParseInt(account[1], 10, 0)
		logindb(nickname, int(id), res)
	} else {
		signupdb(nickname, res)
	}

	if res.Error != packet.Success {
		log.Print("Login End Logout Error")
	}
	if packet.ContainsParamRes(res, "uuid") {
		id = int(res.Param["uuid"].(float64))
	}

	res.Code = packet.Login
	res.Param["nickname"] = nickname
	res.Param["hashTag"] = id
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
	if !packet.ContainsParamReq(res, req, "roomCode") {
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

func (h *Hub) broadcastClientCount() {
	res := packet.NewResponsePacket()
	res.Code = packet.ClientCount
	res.Param["clientCount"] = clientCount
	res.Error = packet.Success
	buf, err := json.Marshal(res)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
	}
	h.broadcast <- buf
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
	case packet.ClientCount:
		c.hub.broadcastClientCount()
	default:
		res.Error = packet.Unknown
	}

	return res
}
