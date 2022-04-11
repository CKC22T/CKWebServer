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

	"github.com/shirou/gopsutil/process"
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
	Err      packet.ErrorCode `json:"err"`
	RoomInfo Room             `json:"roomInfo"`
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

func RoomProfileHandler(w http.ResponseWriter, r *http.Request) {
	var room Room
	json.NewDecoder(r.Body).Decode(&room)
	id, _ := strconv.Atoi(r.URL.Query()["id"][0])
	var pid = dediServers[id].Proc.Pid

	html := "<html>Room id : " + strconv.FormatInt(int64(id), 10) + "<br>"

	pss, _ := process.Processes()
	for _, ps := range pss {
		if ps.Pid == int32(pid) {
			pCpu, _ := ps.CPUPercent()
			pMemInfo, _ := ps.MemoryInfo()
			pMemPer, _ := ps.MemoryPercent()
			pCpuAff, _ := ps.CPUAffinity()
			pConnections, _ := ps.Connections()
			pCwd, _ := ps.Cwd()
			pExe, _ := ps.Exe()
			pGids, _ := ps.Gids()
			pGroups, _ := ps.Groups()
			pIOCounters, _ := ps.IOCounters()
			pIOnice, _ := ps.IOnice()
			//pMemMap, _ := ps.MemoryMaps(true)
			pName, _ := ps.Name()
			pNice, _ := ps.Nice()
			pNumFDs, _ := ps.NumFDs()
			pNumThreads, _ := ps.NumThreads()
			//pPageFaults, _ := ps.PageFaults()
			pParent, _ := ps.Parent()
			pPpid, _ := ps.Ppid()
			pRlimit, _ := ps.Rlimit()
			pRlimitUsage, _ := ps.RlimitUsage(true)
			pStatus, _ := ps.Status()
			pTerminal, _ := ps.Terminal()
			pTgid, _ := ps.Tgid()
			pThreads, _ := ps.Threads()
			pTimes, _ := ps.Times()
			pUids, _ := ps.Uids()
			pUserName, _ := ps.Username()

			html = html + "pid : " + strconv.FormatInt(int64(ps.Pid), 10) + "<br>"
			html = html + "status : " + ps.String() + "<br>"
			html = html + "cpu : " + strconv.FormatFloat(pCpu/16, 'f', 2, 64) + "%<br>"
			for _, aff := range pCpuAff {
				html = html + "cpu aff : " + strconv.FormatInt(int64(aff), 10) + "<br>"
			}
			for idx, conn := range pConnections {
				html = html + "cpu conn [" + strconv.Itoa(idx) + "] : " + conn.Status + "<br>"
			}
			html = html + "memory info : " + pMemInfo.String() + "<br>"
			html = html + "memory per : " + strconv.FormatFloat(float64(pMemPer), 'f', 2, 64) + "%<br>"

			html = html + "Cwd : " + pCwd + "<br>"
			html = html + "Exe : " + pExe + "<br>"
			for idx, guid := range pGids {
				html = html + "Gids[" + strconv.Itoa(idx) + "] : " + strconv.FormatInt(int64(guid), 10) + "<br>"
			}
			for idx, group := range pGroups {
				html = html + "Groups[" + strconv.Itoa(idx) + "] : " + strconv.FormatInt(int64(group), 10) + "<br>"
			}
			html = html + "IOCounters : " + pIOCounters.String() + "<br>"
			html = html + "IOnice : " + strconv.FormatInt(int64(pIOnice), 10) + "<br>"
			// for idx, mem := range *pMemMap {
			// 	html = html + "MemMap[" + strconv.Itoa(idx) + "]" + mem + "<br>"
			// }
			html = html + "Name : " + pName + "<br>"
			html = html + "Nice : " + strconv.FormatInt(int64(pNice), 10) + "<br>"
			html = html + "NumFDs : " + strconv.FormatInt(int64(pNumFDs), 10) + "<br>"
			html = html + "NumThreads : " + strconv.FormatInt(int64(pNumThreads), 10) + "<br>"
			//html = html + "PageFaults ChildMajorFaults : " + strconv.FormatUint(pPageFaults.ChildMajorFaults, 10) + "<br>"
			//html = html + "PageFaults ChildMinorFaults : " + strconv.FormatUint(pPageFaults.ChildMinorFaults, 10) + "<br>"
			//html = html + "PageFaults MajorFaults : " + strconv.FormatUint(pPageFaults.MajorFaults, 10) + "<br>"
			//html = html + "PageFaults MinorFaults : " + strconv.FormatUint(pPageFaults.MinorFaults, 10) + "<br>"

			html = html + "Parent : " + pParent.String() + "<br>"
			html = html + "PPid : " + strconv.FormatInt(int64(pPpid), 10) + "<br>"
			for idx, rlimit := range pRlimit {
				html = html + "Rlimite[" + strconv.Itoa(idx) + "]" + rlimit.String() + "<br>"
			}
			for idx, rlimitUsage := range pRlimitUsage {
				html = html + "RlimiteUsage[" + strconv.Itoa(idx) + "]" + rlimitUsage.String() + "<br>"
			}
			for idx, states := range pStatus {
				html = html + "Rlimite[" + strconv.Itoa(idx) + "]" + states + "<br>"
			}
			html = html + "Terminal : " + pTerminal + "<br>"
			html = html + "Tgid : " + strconv.FormatInt(int64(pTgid), 10) + "<br>"
			for idx, thread := range pThreads {
				html = html + "Rlimite[" + strconv.FormatInt(int64(idx), 10) + "]" + thread.String() + "<br>"
			}
			html = html + "Times : " + pTimes.String() + "<br>"
			for idx, uid := range pUids {
				html = html + "Rlimite[" + strconv.Itoa(idx) + "]" + strconv.FormatInt(int64(uid), 10) + "<br>"
			}
			html = html + "UserName : " + pUserName + "<br>"

			break
		}
	}

	html = html + "</html>"

	w.Write([]byte(html))
}
