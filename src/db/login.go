package db

import (
	"encoding/json"
	"fmt"
	"net/http"
	"packet"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := packet.NewRequestPacket()
	res := packet.NewResponsePacket()

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}
	//내부 로직
	if !packet.ContainsParam(res, req, "nickname") {
		return
	}
	if !packet.ContainsParam(res, req, "id") {
		return
	}
	var nickname string = req.Param["nickname"].(string)
	var id int = int(req.Param["id"].(float64))

	uuid := Login(nickname, id)
	if uuid == 0 {
		res.Error = packet.Unknown
	}

	res.Param["uuid"] = uuid

	response, _ := json.Marshal(res)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}
