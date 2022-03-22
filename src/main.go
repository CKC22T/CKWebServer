package main

import (
	"fmt"
	"lobby"
	"net/http"
	"room"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/login", lobby.LoginHandler)
	http.HandleFunc("/signup", lobby.SignUpHandler)
	http.HandleFunc("/getuserinfo", lobby.GetUserInfoHandler)
	http.HandleFunc("/rooms", room.RoomsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world")
	})
	//http.HandleFunc("/ws", socketHandler)

	http.ListenAndServe(":3000", nil)
}

//SocketHandler
//var upgrader = websocket.upgrader{}

// func SocketHandler(w http.ResponseWriter, r *http.Request) {

// }
