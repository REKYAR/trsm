package main

import (
	"appPrototype1/structures"
	"appPrototype1/utils"
	"crypto/rand"
	"encoding/base64"
	"io"
)

var cookiesInUse = make([]string,0)
var sessions = make(map[string]Session)

type Session struct {
	user structures.User
	lifetime int64
}

func StartSession(uid int64, lifetime int64) (string,error) {
	//zabawa z db

	cookieName := generateId()
	for true {
		_,found := utils.Find(&cookiesInUse,cookieName)
		if !found{
			break
		}
		cookieName=generateId()
	}
	u := structures.User{}
	//fmt.Println(u)
	err:=u.BuildUserFromDB(uid) // tu sie jebie
	//fmt.Println(u)
	if err!=nil {
		return "-1",err
	}
	//86400000 ms =24h
	sessions[cookieName] = Session{user: u, lifetime: lifetime}
	//fmt.Println("current sessions:", sessions)
	return cookieName,nil
}

func generateId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

