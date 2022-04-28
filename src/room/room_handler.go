package room

import (
	"encoding/json"
	"log"
	"net/http"
	"packet"
)

func RoomsHandler(w http.ResponseWriter, r *http.Request) {
	var roomInfo RoomInfo
	var roomInfoRes RoomInfoRes
	json.NewDecoder(r.Body).Decode(&roomInfo)
	switch r.Method {
	case http.MethodGet:
		if roomInfo.Id > 0 {
			roomInfoRes = getRoom(roomInfo.Id)
		} else {
			json.NewEncoder(w).Encode(getRooms())
			return
		}
	case http.MethodPost:
		roomInfoRes = createRoom()
	case http.MethodPut:
		roomInfoRes = updateRoom(roomInfo)
	case http.MethodDelete:
		roomInfoRes = deleteRoom(roomInfo.Id)
	}
	json.NewEncoder(w).Encode(roomInfoRes)
}

func getRooms() map[int]*Room {
	log.Printf("Log[Room] : GetRooms")
	return Rooms
}

func getRoom(roomId int) RoomInfoRes {
	log.Printf("Log[Room] : GetRoom Id:[%d]", roomId)
	var roomInfo RoomInfoRes
	if Rooms[roomId] != nil {
		roomInfo.Err = packet.Unknown
	} else {
		var roomInfo RoomInfoRes
		roomInfo.Err = packet.Success
		roomInfo.RoomInfo = *Rooms[roomId].Info
	}
	return roomInfo
}

func updateRoom(roomInfo RoomInfo) RoomInfoRes {
	log.Printf("Log[Room] : Update Room id:[%d] CurUser:[%d]", roomInfo.Id, roomInfo.CurUser)
	var roomInfoRes RoomInfoRes
	roomInfoRes.Err = packet.Success
	Rooms[roomInfo.Id].Info.CurUser = roomInfo.CurUser
	roomInfoRes.RoomInfo = *Rooms[roomInfo.Id].Info
	if roomInfo.CurUser <= 0 {
		deleteRoom(roomInfo.Id)
	}
	return roomInfoRes
}

func createRoom() RoomInfoRes {
	log.Printf("Log[Room] : Create Room")
	var room Room
	dediProc := DedicatedProcessOnBegin(&room)
	room.Info = &RoomInfo{Id: dediProc.Id, Addr: dediProc.Addr, MaxUser: 4, CurUser: 0}
	Rooms[room.Info.Id] = &room
	var roomInfoRes RoomInfoRes
	roomInfoRes.Err = packet.Success
	roomInfoRes.RoomInfo = *room.Info
	return roomInfoRes
}

func deleteRoom(roomId int) RoomInfoRes {
	log.Printf("Log[Room] : Delete Room id:[%d]", roomId)
	var roomInfoRes RoomInfoRes
	roomInfoRes.Err = packet.Success
	roomInfoRes.RoomInfo = *Rooms[roomId].Info
	DedicatedProcessKill(roomId)
	delete(Rooms, roomId)
	return roomInfoRes
}
