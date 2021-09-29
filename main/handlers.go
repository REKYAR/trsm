package main

import (
	"appPrototype1/dbinteractions"
	"appPrototype1/structures"
	"appPrototype1/utils"
	"crypto/sha256"
	"encoding/base64"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

func feedHandler(writer http.ResponseWriter, request *http.Request)  {
	println(request.Cookies())
	if utils.CheckForCookie(request,"sessionId") {
		ucookie,_ := request.Cookie("sessionId")
		udata := sessions[ucookie.Value]
		html,err :=template.ParseFiles(staticTemplate+"feed.html", staticTemplate+"menu.html")
		utils.ErrCheck(err)
		err=html.Execute(writer,udata.user)
		utils.ErrCheck(err)
	}else {
		html,err :=template.ParseFiles(staticTemplate+"feed.html", staticTemplate+"menu.html")
		utils.ErrCheck(err)
		err=html.Execute(writer,nil)
		utils.ErrCheck(err)
	}
}

func RegisterHandler(writer http.ResponseWriter, request *http.Request){
	println(request.Cookies())
	if request.Method=="GET"{
		if utils.CheckForCookie(request,"sessionId") {
			cmp := Session{}
			ucookie,_ := request.Cookie("sessionId")
			udata := sessions[ucookie.Value]
			if udata!=cmp{
				html,err :=template.ParseFiles(staticTemplate+"feed.html", staticTemplate+"menu.html")
				utils.ErrCheck(err)
				err=html.Execute(writer,udata.user)
				utils.ErrCheck(err)
				return
			}
			html,err :=template.ParseFiles(staticTemplate+"register.html", staticTemplate+"menu.html")
			utils.ErrCheck(err)
			err=html.Execute(writer,nil)
			utils.ErrCheck(err)
			return

		}else {
			html,err :=template.ParseFiles(staticTemplate+"register.html", staticTemplate+"menu.html")
			utils.ErrCheck(err)
			err=html.Execute(writer,nil)
			utils.ErrCheck(err)
			return
		}

	}else if request.Method=="POST"{
		//fmt.Println("POST detected")
		//if logged in then redirect to root

		if utils.CheckForCookie(request,"sessionId") {
			//fmt.Println("inside")
			cmp := Session{}
			ucookie,_ := request.Cookie("sessionId")
			udata := sessions[ucookie.Value]
			if udata!=cmp{
				html,err :=template.ParseFiles(staticTemplate+"feed.html", staticTemplate+"menu.html")
				utils.ErrCheck(err)
				err=html.Execute(writer,udata.user)
				utils.ErrCheck(err)
				return
			}
		}
		login := request.FormValue("login")
		password := request.FormValue("password")
		hash := sha256.New()
		hash.Write([]byte(password))
		//fmt.Println("form data loaded, password hashed")
		newUser := structures.User{
			LoggedIn:      true,
			Id:            0,
			Login:         login,
			TagWhitelist:  "",
			TagBlacklist:  "",
			Rep:           0,
			Password:      base64.URLEncoding.EncodeToString(hash.Sum(nil)),
			Admin:         false,
		}
		//check for user login duplication in db
		var vch = make( chan dbinteractions.UserChannel)
		var vres dbinteractions.UserChannel
		go dbinteractions.GetUserByLogin(login, vch)
		vres = <- vch
		//fmt.Println(vres)
		if vres.UserField.Next() {
			//fmt.Println("user exists, redirecting")
			http.Redirect(writer, request,"/",http.StatusFound)
		}else{
			//fmt.Println("user does not exist")
			//fmt.Println()
			var uch = make(chan dbinteractions.IdChannel)
			var err error
			var ret dbinteractions.IdChannel

			//fmt.Println("vars declared")
			go dbinteractions.AddUser(newUser.Login,newUser.Password,uch)
			//fmt.Println("func on")
			ret=<-uch
			//fmt.Println("user saved to db")
			if ret.Err!=nil {
				utils.ErrCheck(ret.Err)
			}
			//fmt.Println("starting session")
			sessionId,err:= StartSession(ret.Id,86400000)

			utils.ErrCheck(err)
			if sessionId=="-1" {
				log.Fatal("unexpected error has occured")
			}
			//fmt.Println("session started")
			//fmt.Println("almost there")
			c := http.Cookie{Name: "sessionId", Value: sessionId}
			http.SetCookie(writer, &c)
			//fmt.Println("redirecting")
			http.Redirect(writer, request,"/",http.StatusFound)
		}


	}
}

func LogoutHandler(writer http.ResponseWriter, request *http.Request){
	println(request.Cookies())
	if utils.CheckForCookie(request,"sessionId") {
		ucookie,_ := request.Cookie("sessionId")
		delete(sessions, ucookie.Value)
		utils.RemoveFromSlice(&cookiesInUse, ucookie.Name)
		html,err :=template.ParseFiles(staticTemplate+"feed.html", staticTemplate+"menu.html")
		utils.ErrCheck(err)
		err=html.Execute(writer,nil)
		utils.ErrCheck(err)
	}else {
		html,err :=template.ParseFiles(staticTemplate+"feed.html", staticTemplate+"menu.html")
		utils.ErrCheck(err)
		err=html.Execute(writer,nil)
		utils.ErrCheck(err)
	}
}

func LoginHandler(writer http.ResponseWriter, request *http.Request)  {
	//fmt.Println("LOGIN")
	if request.Method=="GET" {
		if utils.CheckForCookie(request,"sessionId") {
			cmp := Session{}
			ucookie,_ := request.Cookie("sessionId")
			udata := sessions[ucookie.Value]
			if udata!=cmp{
				html,err :=template.ParseFiles(staticTemplate+"feed.html", staticTemplate+"menu.html")
				utils.ErrCheck(err)
				err=html.Execute(writer,udata.user)
				utils.ErrCheck(err)
				return
			}
			html,err :=template.ParseFiles(staticTemplate+"login.html", staticTemplate+"menu.html")
			utils.ErrCheck(err)
			err=html.Execute(writer,nil)
			utils.ErrCheck(err)

		}else {
			html,err :=template.ParseFiles(staticTemplate+"login.html", staticTemplate+"menu.html")
			utils.ErrCheck(err)
			err=html.Execute(writer,nil)
			utils.ErrCheck(err)
			return
		}
	}else if request.Method=="POST" {
		if utils.CheckForCookie(request,"sessionId") {
			cmp := Session{}
			ucookie,_ := request.Cookie("sessionId")
			udata := sessions[ucookie.Value]
			if udata!=cmp{
				html,err :=template.ParseFiles(staticTemplate+"feed.html", staticTemplate+"menu.html")
				utils.ErrCheck(err)
				err=html.Execute(writer,udata.user)
				utils.ErrCheck(err)
				return
			}
		}
		login := request.FormValue("login")
		password := request.FormValue("password")
		hash := sha256.New()
		hash.Write([]byte(password))
		hashed:=base64.URLEncoding.EncodeToString(hash.Sum(nil))
		var vch = make( chan dbinteractions.UserChannel)
		var vres dbinteractions.UserChannel

		go dbinteractions.ValidateUserByLogin(login, vch)
		vres = <- vch
		//fmt.Println(vres)
		if vres.UserField.Next() {
			var id int64
			var passwordFromDatabase string
			//fmt.Println("user exists")
			//fmt.Println()
			var err error
			vres.UserField.Scan(&id,&passwordFromDatabase)
			if passwordFromDatabase!=hashed {
				html,err :=template.ParseFiles(staticTemplate+"login.html", staticTemplate+"menu.html")
				utils.ErrCheck(err)
				err=html.Execute(writer,nil)
				utils.ErrCheck(err)
				return
			}

			//fmt.Println("starting session")
			sessionId,err:= StartSession(id,86400000)

			utils.ErrCheck(err)
			if sessionId=="-1" {
				log.Fatal("unexpected error has occured")
			}
			//fmt.Println("session started")
			//fmt.Println("almost there")
			c := http.Cookie{Name: "sessionId", Value: sessionId}
			http.SetCookie(writer, &c)
			//fmt.Println("redirecting")
			http.Redirect(writer, request,"/",http.StatusFound)
		}else{
			//fmt.Println("user does not exist, redirecting")
			http.Redirect(writer, request,"/register",http.StatusFound)

		}
	}
}

func DeleteAccHandler(writer http.ResponseWriter, request *http.Request)  {
	if utils.CheckForCookie(request,"sessionId") {
		ucookie,_ := request.Cookie("sessionId")
		udata := sessions[ucookie.Value]
		sch := make(chan dbinteractions.StringChannel)
		ech := make(chan error)
		go dbinteractions.Delpfp(int64(udata.user.Id),sch)
		data:=<-sch
		utils.ErrCheck(data.Err)
		if data.Str!="-1" {
			err:=os.Remove("pfp_storage/"+data.Str)
			utils.ErrCheck(err)
		}
		delete(sessions, ucookie.Value)
		utils.RemoveFromSlice(&cookiesInUse, ucookie.Name)
		go dbinteractions.DeleteUserById(int64(udata.user.Id),ech)
		//fmt.Println("session closed")
		err:=<-ech//
		utils.ErrCheck(err)
		http.Redirect(writer, request,"/register",http.StatusFound)
	}else {
		http.Redirect(writer, request,"/register",http.StatusFound)
	}
}


func ProfileHandler(writer http.ResponseWriter, request *http.Request)  {
	cusr := structures.User{}
	if utils.CheckForCookie(request,"sessionId") {
		cmp := Session{}
		ucookie,_ := request.Cookie("sessionId")
		udata := sessions[ucookie.Value]
		if udata!=cmp{
			cusr = udata.user
		}

	}

	suid :=  mux.Vars(request)["id"]
	uid, err := strconv.ParseInt(suid, 10, 64)
	utils.ErrCheck(err)
	//spr czy user ma pfp
	var searchedFor structures.User
	err=searchedFor.BuildUserFromDB(uid)
	utils.ErrCheck(err)
	//path, err := os.Getwd()
	//utils.ErrCheck(err)
	//to jak nie ma pfp
	//default pfp
	var path string
	sch := make(chan dbinteractions.StringChannel)
	go dbinteractions.GetPfp(int64(cusr.Id),sch)
	rsch := <-sch
	utils.ErrCheck(rsch.Err)
	if rsch.Str=="-1" {
		path = "defaut_pfp.png"
	}else {
		path = rsch.Str
	}
	//{{.Inspected.PfpPath}}
	//fmt.Println(path)
	upd := structures.UserProfileData{
		User:      cusr,
		Inspected: structures.Inspected{
			Id:      uid,
			Login:   searchedFor.Login,
			Rep:     searchedFor.Rep,
			PfpPath: path,
		},
	}
	html,err :=template.ParseFiles(staticTemplate+"profile.html", staticTemplate+"menu.html")
	utils.ErrCheck(err)
	err=html.Execute(writer,upd)
	utils.ErrCheck(err)

}

func ProfileEditHandler(writer http.ResponseWriter, request *http.Request)  {
	if utils.CheckForCookie(request,"sessionId") {
		cmp := Session{}
		ucookie,_ := request.Cookie("sessionId")
		udata := sessions[ucookie.Value]
		if udata!=cmp{
			html,err :=template.ParseFiles(staticTemplate+"edit_profile.html", staticTemplate+"menu.html")
			utils.ErrCheck(err)
			err=html.Execute(writer,udata.user)
			utils.ErrCheck(err)
			return
		}
		html,err :=template.ParseFiles(staticTemplate+"login.html", staticTemplate+"menu.html")
		utils.ErrCheck(err)
		err=html.Execute(writer,nil)
		utils.ErrCheck(err)
		return
	}else {
		html,err :=template.ParseFiles(staticTemplate+"login.html", staticTemplate+"menu.html")
		utils.ErrCheck(err)
		err=html.Execute(writer,nil)
		utils.ErrCheck(err)
		return
	}
}

