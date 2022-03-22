package lobby

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"packet"
)

type SignInReq struct {
	Id string `json:"id"`
	Pw string `json:"pw"`
}
type SignInRes struct {
	Err  int `json:"err"`
	Uuid int `json:"uuid"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	requestData := new(SignInReq)
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}
	//내부 로직
	responseData := new(SignInRes)
	responseData.Err = packet.Success

	uuid := db.Login(requestData.Id, requestData.Pw)
	if uuid == 0 {
		responseData.Err = packet.Unknown
	}

	responseData.Uuid = uuid

	response, _ := json.Marshal(responseData)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}
