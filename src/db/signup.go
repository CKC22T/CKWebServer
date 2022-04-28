package db

import (
	"encoding/json"
	"fmt"
	"net/http"
	"packet"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	req := packet.NewRequestPacket()
	res := packet.NewResponsePacket()

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}

	if !packet.ContainsParam(res, req, "nickname") {
		return
	}
	var nickname string = req.Param["nickname"].(string)

	result := SignUp(nickname)

	if result {
		res.Error = packet.Success
	} else {
		res.Error = packet.Unknown
	}

	response, _ := json.Marshal(res)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}
