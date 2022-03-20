package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

const (
	Success = iota
	Unknown
)

//var db *sql.DB
var dbroot string = "root:gomjun423009@tcp(gomjun.asuscomm.com:3306)/olympus"

func main() {
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/signup", SignUpHandler)
	http.HandleFunc("/createroom", CreateRoomHandler)
	http.HandleFunc("/getuserinfo", GetUserInfoHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world")
	})
	//http.HandleFunc("/ws", socketHandler)

	http.ListenAndServe(":3000", nil)
}

//db func
func SignUp(id string, pw string, nick string) bool {
	db, err := sql.Open("mysql", dbroot)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	result, err := db.Exec("INSERT INTO user(create_time, update_time, account_id, account_pw, nickname) VALUES(now(), now(), '" + id + "', '" + pw + "', '" + nick + "')")
	if err != nil {
		log.Fatal(err)
		return false
	}
	result.RowsAffected()

	return true
}

func Login(id string, pw string) int {
	db, err := sql.Open("mysql", dbroot)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	var uuid int
	err = db.QueryRow("SELECT id FROM user WHERE account_id = '" + id + "' AND account_pw = '" + pw + "'").Scan(&uuid)
	if err != nil {
		return 0
	}

	return uuid
}

type DTOUserInfo struct {
	Nick string
}

func GetUserData(uuid int) *DTOUserInfo {
	db, err := sql.Open("mysql", dbroot)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	var nickname string
	err = db.QueryRow("SELECT nickname FROM user WHERE id = '" + strconv.Itoa(uuid) + "'").Scan(&nickname)

	var userInfo = new(DTOUserInfo)
	userInfo.Nick = nickname

	if err != nil {
		return userInfo
	}

	return userInfo
}

//login
type SignInReq struct {
	Id string
	Pw string
}
type SignInRes struct {
	Err  int `json:"err"`
	Uuid int
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
	responseData.Err = Success

	uuid := Login(requestData.Id, requestData.Pw)
	if uuid == 0 {
		responseData.Err = Unknown
	}

	responseData.Uuid = uuid

	response, _ := json.Marshal(responseData)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}

//signup
type SignUpReq struct {
	Id   string
	Pw   string
	Nick string
}

type SignUpRes struct {
	Err int
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	requestData := new(SignUpReq)
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}

	result := SignUp(requestData.Id, requestData.Pw, requestData.Nick)

	responseData := new(SignUpRes)
	if result {
		responseData.Err = Success
	} else {
		responseData.Err = Unknown
	}

	response, _ := json.Marshal(responseData)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}

//GetUserInfo
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

	data := GetUserData(requestData.Uuid)

	responseData := new(UserInfoRes)
	responseData.Nick = data.Nick
	responseData.Err = Success

	response, _ := json.Marshal(responseData)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(response))
}

//POST /rooms
type CreateRoomReq struct {
	Name string
}

type CreateRoomRes struct {
	Err int `json:"err"`
	Id  int
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	requestData := new(CreateRoomReq)
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}
}

//GET /rooms

//PUT /rooms

//DELETE /rooms

//SocketHandler
//var upgrader = websocket.upgrader{}

// func SocketHandler(w http.ResponseWriter, r *http.Request) {

// }
