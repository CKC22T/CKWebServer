package room

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Room struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	MaxUser int    `json:"maxUser"`
	CurUser int    `json:"CurUser"`
}

var roomAtoi int = 0
var rooms = map[int]*Room{}

func RoomsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var room Room
		json.NewDecoder(r.Body).Decode(&room)
		fmt.Println(room.Id)
		json.NewEncoder(w).Encode(rooms)
	case http.MethodPost:
		var room Room
		json.NewDecoder(r.Body).Decode(&room)
		room.Id = roomAtoi
		rooms[room.Id] = &room
		roomAtoi = roomAtoi + 1
		json.NewEncoder(w).Encode(room)
	case http.MethodPut:
		var room Room
		json.NewDecoder(r.Body).Decode(&room)
		rooms[room.Id] = &room
		json.NewEncoder(w).Encode(room)
	case http.MethodDelete:
		var room Room
		json.NewDecoder(r.Body).Decode(&room)
		delete(rooms, room.Id)
		json.NewEncoder(w).Encode(room)
	}
}
