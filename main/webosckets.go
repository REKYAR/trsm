package main

import (
	"appPrototype1/dbinteractions"
	"appPrototype1/utils"
	"crypto/sha256"
	"encoding/base64"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
func PasswordEditHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.ErrCheck(err)
	_, p, err := conn.ReadMessage() //mode 1
	utils.ErrCheck(err)
	val :=string(p)
	if utils.CheckForCookie(r,"sessionId") {
		cmp := Session{}
		ucookie,_ := r.Cookie("sessionId")
		udata := sessions[ucookie.Value]
		if udata!=cmp{
			hash := sha256.New()
			hash.Write([]byte(val))
			newPassword:=base64.URLEncoding.EncodeToString(hash.Sum(nil))
			tmpSession:=sessions[ucookie.Value]
			tmpSession.user.Password=newPassword
			sessions[ucookie.Value]=tmpSession
			ech := make(chan error)
			go dbinteractions.ChangePassword(int64(sessions[ucookie.Value].user.Id),newPassword,ech)
			err= <-ech
			utils.ErrCheck(err)
		}
	}
}

func PfpEditHandler(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.ErrCheck(err)
	if utils.CheckForCookie(r,"sessionId") {
		cmp := Session{}
		ucookie,_ := r.Cookie("sessionId")
		udata := sessions[ucookie.Value]
		if udata!=cmp{
			_, imagePassed, err := conn.ReadMessage() //mode 2
			utils.ErrCheck(err)
			pfpname:=utils.GenerateRandomName(16)
			path,err := os.Getwd()
			utils.ErrCheck(err)
			for true {
				if _, err := os.Stat(path+"/pfpstorage"+pfpname); os.IsNotExist(err) {
					break
				}
			}
			err = ioutil.WriteFile("pfp_storage/"+pfpname, imagePassed,0777)
			utils.ErrCheck(err)
			sch := make(chan dbinteractions.StringChannel)
			go dbinteractions.AddPfp(int64(udata.user.Id),pfpname,sch)
			data := <- sch
			//fmt.Println(data.Str)
			utils.ErrCheck(data.Err)
			if data.Str!="-1" {
				err=os.Remove("pfp_storage/"+data.Str)
				utils.ErrCheck(data.Err)
			}

		}
	}


}