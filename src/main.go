package main

import (
	"fmt"
	"lobby"
	"net/http"
	_ "net/http/pprof"
	"profile"
	"room"
	"websocket"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	room.SetIp()
	hub := websocket.NewHub()
	go hub.Run()

	http.HandleFunc("/rooms", room.RoomsHandler)
	http.HandleFunc("/user", lobby.GetUserInfoHandler)
	http.HandleFunc("/open", room.DedicatedProcessOnEnd)

	http.HandleFunc("/profile", profile.ProfileHandler)
	http.HandleFunc("/room/profile", profile.RoomProfileHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServerWs(hub, w, r)
	})

	http.ListenAndServe(":3000", nil)
}
