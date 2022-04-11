package lobby

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"packet"
)

type SignUpReq struct {
	Id   string `json:"id"`
	Pw   string `json:"pw"`
	Nick string `json:"nick"`
}

type SignUpRes struct {
	Err packet.ErrorCode `json:"err"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	requestData := new(SignUpReq)
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}

	result := db.SignUp(requestData.Id, requestData.Pw, requestData.Nick)

	responseData := new(SignUpRes)
	if result {
		responseData.Err = packet.Success
	} else {
		responseData.Err = packet.Success
	}

	response, _ := json.Marshal(responseData)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}
