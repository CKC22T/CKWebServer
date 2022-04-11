package main

import (
	"fmt"
	"lobby"
	"net/http"
	_ "net/http/pprof"
	"os"
	"profile"
	"room"
	"websocket"

	_ "github.com/go-sql-driver/mysql"
)

type Address struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type DediProc struct {
	Proc *os.Process
	Id   int
	IsOn chan bool
	Addr Address
}

type ResponsePacket struct {
	Err   int                 `json:"err"`
	Param map[string]struct{} `json:"param"`
}

var dediServers = map[int]*DediProc{}
var dediCodeCount = 0
var dediInitPort = 16000

func main() {
	room.SetIp()
	hub := websocket.NewHub()
	go hub.Run()

	http.HandleFunc("/login", lobby.LoginHandler)
	http.HandleFunc("/signup", lobby.SignUpHandler)
	http.HandleFunc("/getuserinfo", lobby.GetUserInfoHandler)
	http.HandleFunc("/rooms", room.RoomsHandler)
	http.HandleFunc("/user", lobby.GetUserInfoHandler)
	//http.HandleFunc("/run", ProcessHandler)
	http.HandleFunc("/open", room.DedicatedProcessOnEnd)

	http.HandleFunc("/profile", profile.ProfileHandler)
	http.HandleFunc("/room/profile", room.RoomProfileHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServerWs(hub, w, r)
	})

	http.ListenAndServe(":3000", nil)
}
