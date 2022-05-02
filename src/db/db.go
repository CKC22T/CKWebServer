package db

import (
	"database/sql"
	"log"
	"strconv"
)

var dbroot string = "root:gomjun423009@tcp(gomjun.asuscomm.com:3306)/olympus"

func SignUp(nickname string) int {
	db, err := sql.Open("mysql", dbroot)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO user(nickname) VALUES(?)", nickname)
	id, _ := result.LastInsertId()
	if err != nil {
		log.Print(err.Error())
		return 0
	}
	result.RowsAffected()

	return int(id)
}

func Login(nickname string, id int) int {
	if id <= 0 {
		return 0
	}

	db, err := sql.Open("mysql", dbroot)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var uuid int
	err = db.QueryRow("SELECT id FROM `user` WHERE nickname = ? AND id = ?", nickname, id).Scan(&uuid)
	if err != nil {
		return 0
	}
	db.QueryRow("UPDATE `user` SET access_count = access_count + 1, update_time = NOW() WHERE id = ?", id)

	return uuid
}

func ChangeNickname(nickname string, id int) bool {
	if id <= 0 {
		return false
	}
	id_str := strconv.FormatInt(int64(id), 10)

	db, err := sql.Open("mysql", dbroot)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.QueryRow("UPDATE `user` SET nickname = " + nickname + ", update_time = NOW() WHERE id = " + id_str)

	return true
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

func Log(logType string, targetType string, targetCode int, logJsonData string) bool {
	db, err := sql.Open("mysql", dbroot)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO log VALUES(now(), ?, ?, ?, ?)", logType, targetType, targetCode, logJsonData)
	if err != nil {
		log.Fatal(err)
		return false
	}
	result.RowsAffected()

	return true
}

func JoinTime(id int, time float64) bool {
	db, err := sql.Open("mysql", dbroot)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO join_time VALUES(?, now(), ?)", id, time)
	if err != nil {
		log.Fatal(err)
		return false
	}
	result.RowsAffected()

	return true
}
