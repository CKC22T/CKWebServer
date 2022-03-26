package main

import (
	"encoding/json"
	"fmt"
	"lobby"
	"log"
	"net/http"
	"os"
	"os/exec"
	"room"
	"strconv"
	"syscall"

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
	http.HandleFunc("/login", lobby.LoginHandler)
	http.HandleFunc("/signup", lobby.SignUpHandler)
	http.HandleFunc("/getuserinfo", lobby.GetUserInfoHandler)
	http.HandleFunc("/rooms", room.RoomsHandler)
	http.HandleFunc("/run", ProcessHandler)
	http.HandleFunc("/open", DedicatedProcessOnHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world")
	})
	//http.HandleFunc("/ws", socketHandler)

	http.ListenAndServe(":3000", nil)
}

func ProcessHandler(w http.ResponseWriter, r *http.Request) {
	process := exec.Command("../../Build/CKDedi/CKC2022.exe", strconv.Itoa(dediCodeCount))
	process.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 16, NoInheritHandles: true}

	// process.SysProcAttr.CreationFlags = 16
	// process.SysProcAttr.HideWindow = false
	// process.SysProcAttr.NoInheritHandles = true

	err := process.Start()
	if err != nil {
		log.Println(err)
	}
	dediProc := &DediProc{Proc: process.Process, Id: dediCodeCount, Addr: Address{"0.0.0.0", dediInitPort + dediCodeCount}}
	dediServers[dediCodeCount] = dediProc
	dediCodeCount = dediCodeCount + 1

	dediProc.IsOn = make(chan bool)
	isOn := <-dediProc.IsOn
	if isOn {
		log.Printf("[Dedi] Dedicated Server On : ID[%v] IP[%v] Port[%v]", dediProc.Id, dediProc.Addr.Ip, dediProc.Addr.Port)
	} else {
		log.Println("[Error] Dedicated Server Not Open")
	}

	json.NewEncoder(w).Encode(dediProc.Addr)
}

func DedicatedProcessOnHandler(w http.ResponseWriter, r *http.Request) {
	var addr Address
	json.NewDecoder(r.Body).Decode(&addr)
	var procId = addr.Port - dediInitPort
	dediServers[procId].Addr.Ip = addr.Ip
	dediServers[procId].IsOn <- true
	json.NewEncoder(w).Encode(true)
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
