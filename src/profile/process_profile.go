package profile

import (
	"encoding/json"
	"net/http"
	"room"
	"strconv"

	"github.com/shirou/gopsutil/process"
)

func RoomProfileHandler(w http.ResponseWriter, r *http.Request) {
	var req room.Room
	json.NewDecoder(r.Body).Decode(&req)
	id, _ := strconv.Atoi(r.URL.Query()["id"][0])
	var pid = room.Rooms[id].Process.Proc.Process.Pid

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
			html = html + "cpu : " + strconv.FormatFloat(pCpu, 'f', 2, 64) + "%<br>"
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

	html = html + "<br>"
	html = html + "<br>"

	html = html + "</html>"

	w.Write([]byte(html))
}
