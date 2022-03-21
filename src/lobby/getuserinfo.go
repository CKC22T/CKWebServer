package lobby

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"packet"
)

type UserInfoReq struct {
	Uuid int `json:"uuid"`
}

type UserInfoRes struct {
	Err  int    `json:"err"`
	Nick string `json:"nick"`
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	requestData := new(UserInfoReq)
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}

	data := db.GetUserData(requestData.Uuid)

	responseData := new(UserInfoRes)
	responseData.Nick = data.Nick
	responseData.Err = packet.Success

	response, _ := json.Marshal(responseData)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}
