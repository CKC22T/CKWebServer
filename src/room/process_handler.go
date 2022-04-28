package room

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"packet"
	"strconv"
	"syscall"
)

func DedicatedProcessOnBegin(room *Room) *DediProc {
	process := exec.Command("../../Build/ServerBuild/CKC2022.exe", strconv.Itoa(packet.START_BY_WEB), strconv.Itoa(dediCodeCount))
	//process := exec.Command("../../Build/ServerBuild/CKC2022.exe", strconv.Itoa(dediCodeCount))
	process.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 16, NoInheritHandles: true}

	err := process.Start()
	if err != nil {
		log.Println(err)
	}
	dediProc := &DediProc{Proc: process, Id: dediCodeCount, Addr: Address{dediServerIp, dediInitPort + dediCodeCount}}
	room.Info = &RoomInfo{Id: dediProc.Id, Addr: dediProc.Addr, MaxUser: 4, CurUser: 0}
	room.Process = dediProc
	dediCodeCount = dediCodeCount + 1

	Rooms[room.Info.Id] = room
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
	if Rooms[procId] != nil {
		Rooms[procId].Process.IsOn <- true
	} else {
		var room Room
		room.Info = &RoomInfo{Id: procId + 1000, Addr: addr, MaxUser: 4, CurUser: 0}
		Rooms[room.Info.Id] = &room
		result = false
	}
	json.NewEncoder(w).Encode(result)
}

func DedicatedProcessKill(id int) {
	Rooms[id].Process.Proc.Process.Kill()
	delete(Rooms, id)
}
