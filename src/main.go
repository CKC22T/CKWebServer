package main

import (
	"fmt"
	"lobby"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"profile"
	"room"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
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
	//http.HandleFunc("/ws", socketHandler)

	http.ListenAndServe(":3000", nil)
}

//SocketHandler
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrader.Upgrader: %v", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		fmt.Println(string(p))

		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Printf("conn.WriteMessage: %v", err)
			return
		}
	}
}
