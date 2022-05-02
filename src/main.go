package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"lobby"
	"net/http"
	_ "net/http/pprof"
	"profile"
	"room"
	"strconv"
	"time"
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

	http.HandleFunc("/download", ForceDownload)

	http.ListenAndServe(":3000", nil)
}

func ForceDownload(w http.ResponseWriter, r *http.Request) {
	file := "ClientBuild.zip"
	downloadBytes, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println(err)
	}

	// set the default MIME type to send
	mime := http.DetectContentType(downloadBytes)

	fileSize := len(string(downloadBytes))

	// Generate the server headers
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename="+file+"")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	//b := bytes.NewBuffer(downloadBytes)
	//if _, err := b.WriteTo(w); err != nil {
	//              fmt.Fprintf(w, "%s", err)
	//      }

	// force it down the client's.....
	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(downloadBytes))
}
