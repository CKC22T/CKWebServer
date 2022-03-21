package db

import (
	"database/sql"
	"log"
	"strconv"
)

var dbroot string = "root:gomjun423009@tcp(gomjun.asuscomm.com:3306)/olympus"

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
