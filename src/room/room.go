package room

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"packet"
	"strconv"
	"syscall"
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
var dediServers = map[int]*DediProc{}
var dediCodeCount = 1
var dediInitPort = 16000

type Room struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	MaxUser int     `json:"maxUser"`
	CurUser int     `json:"curUser"`
	Addr    Address `json:"addr"`
}

type RoomInfoRes struct {
	Err      int  `json:"err"`
	RoomInfo Room `json:"roomInfo"`
}

var rooms = map[int]*Room{}

func RoomsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var room Room
		json.NewDecoder(r.Body).Decode(&room)
		fmt.Println(room.Id)
		if room.Id > 0 {
			if rooms[room.Id] != nil {
				var roomInfo RoomInfoRes
				roomInfo.Err = packet.Unknown
				json.NewEncoder(w).Encode(roomInfo)
			} else {
				var roomInfo RoomInfoRes
				roomInfo.Err = packet.Success
				roomInfo.RoomInfo = *rooms[room.Id]
				json.NewEncoder(w).Encode(roomInfo)
			}
		} else {
			json.NewEncoder(w).Encode(rooms)
		}
	case http.MethodPost:
		var room Room
		json.NewDecoder(r.Body).Decode(&room)
		dediProc := DedicatedProcessOnBegin()
		room.Id = dediProc.Id
		room.Addr = dediProc.Addr
		room.MaxUser = 4
		room.CurUser = 0
		rooms[room.Id] = &room
		var roomInfo RoomInfoRes
		roomInfo.Err = packet.Success
		roomInfo.RoomInfo = room
		json.NewEncoder(w).Encode(roomInfo)
	case http.MethodPut:
		var room Room
		json.NewDecoder(r.Body).Decode(&room)
		rooms[room.Id].CurUser = room.CurUser
		var roomInfo RoomInfoRes
		roomInfo.Err = packet.Success
		roomInfo.RoomInfo = *rooms[room.Id]
		json.NewEncoder(w).Encode(roomInfo)
	case http.MethodDelete:
		var room Room
		json.NewDecoder(r.Body).Decode(&room)
		DedicatedProcessKill(room.Id)
		delete(rooms, room.Id)
		json.NewEncoder(w).Encode(room)
	}
}

func DedicatedProcessOnBegin() *DediProc {
	process := exec.Command("../../Build/ServerBuild/CKC2022.exe", strconv.Itoa(dediCodeCount))
	process.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 16, NoInheritHandles: true}

	err := process.Start()
	if err != nil {
		log.Println(err)
	}
	dediProc := &DediProc{Proc: process.Process, Id: dediCodeCount, Addr: Address{dediServerIp, dediInitPort + dediCodeCount}}
	dediServers[dediCodeCount] = dediProc
	dediCodeCount = dediCodeCount + 1

	dediProc.IsOn = make(chan bool)
	isOn := <-dediProc.IsOn
	if isOn {
		log.Printf("[Dedi] Dedicated Server On : ID[%v] IP[%v] Port[%v]", dediProc.Id, dediProc.Addr.Ip, dediProc.Addr.Port)
	} else {
		log.Println("[Error] Dedicated Server Not Open")
	}

	return dediProc
}

func DedicatedProcessOnEnd(w http.ResponseWriter, r *http.Request) {
	var addr Address
	json.NewDecoder(r.Body).Decode(&addr)
	var procId = addr.Port - dediInitPort
	var result bool = true
	if dediServers[procId] != nil {
		dediServers[procId].IsOn <- true
	} else {
		var room Room
		room.Id = 100 + procId
		room.Addr = addr
		room.MaxUser = 4
		room.CurUser = 0
		rooms[room.Id] = &room
		result = false
	}
	json.NewEncoder(w).Encode(result)
}

func DedicatedProcessKill(id int) {
	dediServers[id].Proc.Kill()
}
