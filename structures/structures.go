package structures

import (
	"appPrototype1/dbinteractions"
	"fmt"
	"time"
)

type User struct {
	LoggedIn bool
	Id int
	Login string
	TagWhitelist string
	TagBlacklist string
	Rep int
	Password string
	Admin bool
}

type UserReciver struct{
	User_id int
	User_login string
	User_password string
	Created_on time.Time
	Last_login time.Time
	Tag_whitelist string
	Tag_blacklist string
	Admin bool
	Rep int
}

type Inspected struct {
	Id int64
	Login string
	Rep int
	PfpPath string
}

type UserProfileData struct {
	User User
	Inspected Inspected
}

func (u *User) BuildUserFromDB(uid int64) error  {
	var res dbinteractions.UserChannel
	var uch  = make(chan dbinteractions.UserChannel)
	var userRecived UserReciver
	go dbinteractions.GetUserById(uid,uch)
	res= <- uch
	if res.Err!=nil {
		return res.Err
	}
	res.UserField.Next()
	err:=res.UserField.Scan(&userRecived.User_id, &userRecived.User_login, &userRecived.User_password, &userRecived.Created_on,
		&userRecived.Last_login, &userRecived.Tag_whitelist, &userRecived.Tag_whitelist, &userRecived.Admin, &userRecived.Rep)
	if err!=nil {
		return err
	}
	u.TagWhitelist = userRecived.Tag_whitelist
	u.Admin = userRecived.Admin
	u.TagBlacklist = userRecived.Tag_blacklist
	u.Login = userRecived.User_login
	u.Id = userRecived.User_id
	u.LoggedIn = true
	u.Rep = userRecived.Rep
	u.Password=userRecived.User_password
	fmt.Println(u)
	return nil
}