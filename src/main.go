package main

import (
	"encoding/json"
	"fmt"
	"lobby"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/login", lobby.LoginHandler)
	http.HandleFunc("/signup", lobby.SignUpHandler)
	http.HandleFunc("/getuserinfo", lobby.GetUserInfoHandler)
	http.HandleFunc("/createroom", CreateRoomHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world")
	})
	//http.HandleFunc("/ws", socketHandler)

	http.ListenAndServe(":3000", nil)
}

//POST /rooms
type CreateRoomReq struct {
	Name string
}

type CreateRoomRes struct {
	Err int `json:"err"`
	Id  int
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	requestData := new(CreateRoomReq)
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}
}

//GET /rooms

//PUT /rooms

//DELETE /rooms

//SocketHandler
//var upgrader = websocket.upgrader{}

// func SocketHandler(w http.ResponseWriter, r *http.Request) {

// }
