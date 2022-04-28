package room

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"packet"
)

type Address struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type DediProc struct {
	Proc *exec.Cmd
	Id   int
	IsOn chan bool
	Addr Address
}

type IP struct {
	Query string
}

func SetIp() {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ip IP
	json.Unmarshal(body, &ip)

	dediServerIp = ip.Query
}

var dediServerIp string

//var dediServers = map[int]*DediProc{}
var dediCodeCount = 1
var dediInitPort = 50000

type Room struct {
	Info    *RoomInfo
	Process *DediProc
}

type RoomInfo struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	MaxUser int     `json:"maxUser"`
	CurUser int     `json:"curUser"`
	Addr    Address `json:"addr"`
}

type RoomInfoRes struct {
	Err      packet.ErrorCode `json:"err"`
	RoomInfo RoomInfo         `json:"roomInfo"`
}

var Rooms = map[int]*Room{}
