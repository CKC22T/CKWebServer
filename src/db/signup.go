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

	if !packet.ContainsParamReq(res, req, "nickname") {
		return
	}
	var nickname string = req.Param["nickname"].(string)

	uuid := SignUp(nickname)

	if uuid == 0 {
		res.Error = packet.Unknown
	} else {
		res.Error = packet.Success
	}
	res.Param["uuid"] = uuid

	response, _ := json.Marshal(res)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}
