package db

import (
	"encoding/json"
	"fmt"
	"net/http"
	"packet"
)

type LogDTO struct {
	LogType     string `json:"logType"`
	TargetType  string `json:"targetType"`
	TargetCode  int    `json:"targetCode"`
	LogJsonData string `json:"logJsonData"`
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	var req LogDTO
	res := packet.NewResponsePacket()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}
	//내부 로직
	logType := req.LogType
	targetType := req.TargetType
	targetCode := req.TargetCode
	logJsonData := req.LogJsonData

	result := Log(logType, targetType, targetCode, logJsonData)

	if !result {
		res.Error = packet.Unknown
	}

	response, _ := json.Marshal(res)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}

type JoinTimeDTO struct {
	Id   int     `json:"id"`
	Time float64 `json:"time"`
}

func JoinTimeHandler(w http.ResponseWriter, r *http.Request) {
	var req JoinTimeDTO
	res := packet.NewResponsePacket()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}
	//내부 로직
	id := req.Id
	time := req.Time

	result := JoinTime(id, time)

	if !result {
		res.Error = packet.Unknown
	}

	response, _ := json.Marshal(res)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}
