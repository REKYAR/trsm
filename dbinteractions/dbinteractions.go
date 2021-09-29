package dbinteractions

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)



func openDb()( *sql.DB,  error){
	db, err := sql.Open("postgres", os.Getenv("HEROKU_POSTGRESQL_PURPLE_URL"))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, nil
}

func closeDb(db *sql.DB) error  {
	err:=db.Close()
	return err
}

func closeStatement(stmt *sql.Stmt)  error{
	err:=stmt.Close()
	return err
}

func AddUser(Login string, Password string, uch chan IdChannel) {
	db,err := openDb()
	if err!=nil {
		uch<-IdChannel{Id: -1, Err: err}
	}
	stmt, err := db.Prepare("INSERT INTO users (user_login,user_password,created_on,last_login,tag_whitelist,tag_blacklist,admin,rep ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING user_id")
	if err!=nil {uch<-IdChannel{Id: -1, Err: err}}
	var res int64
	err = stmt.QueryRow(Login, Password, time.Now(),time.Now(),"","",false,0).Scan(&res)
	if err!=nil {
		uch<-IdChannel{Id: -1, Err: err}
	}
	err = closeStatement(stmt)
	if err!=nil {uch<-IdChannel{Id: -1, Err: err}}
	err = closeDb(db)
	if err!=nil {
		//fmt.Println(err)
		uch<-IdChannel{Id: -1, Err: err}
	}
	uch<-IdChannel{Id: res, Err: nil}
}

func GetUserById(id int64, uch chan UserChannel)  {
	db,err := openDb()
	if err!=nil {uch<-UserChannel{UserField: nil, Err: err}}
	res,err := db.Query("SELECT * FROM users WHERE user_id=$1",id)
	if err!=nil {uch<-UserChannel{UserField: nil, Err: err}}
	err = closeDb(db)
	if err!=nil {uch<-UserChannel{UserField: nil, Err: err}}
	uch<-UserChannel{UserField: res, Err: nil}
}

func GetUserByLogin(login string, uch chan UserChannel)  {
	db,err := openDb()
	if err!=nil {uch<-UserChannel{UserField: nil, Err: err}}
	res,err := db.Query("SELECT * FROM users WHERE user_login=$1",login)
	if err!=nil {uch<-UserChannel{UserField: nil, Err: err}}
	err = closeDb(db)
	if err!=nil {uch<-UserChannel{UserField: nil, Err: err}}
	uch<-UserChannel{UserField: res, Err: nil}
}

func ValidateUserByLogin(login string, uch chan UserChannel)  {
	db,err := openDb()
	if err!=nil {uch<-UserChannel{UserField: nil, Err: err}}
	res,err := db.Query("SELECT user_id,user_password FROM users WHERE user_login=$1",login)
	if err!=nil {uch<-UserChannel{UserField: nil, Err: err}}
	err = closeDb(db)
	if err!=nil {uch<-UserChannel{UserField: nil, Err: err}}
	uch<-UserChannel{UserField: res, Err: nil}
}

func DeleteUserById(id int64, ech chan error)  {
	db,err := openDb()
	if err!=nil {ech<-err}
	_,err = db.Exec("DELETE FROM users WHERE user_id=$1",id)
	if err!=nil {ech<-err}
	err = closeDb(db)
	if err!=nil {ech<-err}
	ech<-nil
}

func ChangePassword(id int64,newPassword string, ech chan error)  {
	db,err := openDb()
	if err!=nil {ech<-err}
	_,err = db.Exec("UPDATE users SET user_password=$1 WHERE user_id=$2;",newPassword,id)
	if err!=nil {ech<-err}
	err = closeDb(db)
	if err!=nil {ech<-err}
	ech<-nil
}

//removal of existing pic not working, first time upoload is ok
func AddPfp(id int64,filename string, sch chan StringChannel)  {
	db,err := openDb()
	if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	res,err := db.Query("SELECT pfp FROM users_pfps WHERE user_id=$1",id)
	if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	var tmp string
	if res.Next() {
		//fmt.Println("user has a pfp, removing...")
		err=res.Scan(&tmp) //why tho?
		if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
		_,err = db.Exec("DELETE FROM users_pfps WHERE user_id=$1;",id)
		if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	}else {
		tmp ="-1"
	}
	//fmt.Println("inserting new pfp...")
	_,err = db.Exec("INSERT INTO users_pfps (pfp,user_id) VALUES ($1,$2);",filename,id)
	if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	err = closeDb(db)
	if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	sch<-StringChannel{Err: nil, Str: tmp}
}

//to be checked( ok)
func GetPfp(id int64, sch chan StringChannel){
	db,err := openDb()
	if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	res,err := db.Query("SELECT pfp FROM users_pfps WHERE user_id=$1",id)
	if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	var tmp string
	if res.Next() {
		err=res.Scan(&tmp)
		if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	}else {
		tmp ="-1"
	}
	err = closeDb(db)
	if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	sch<-StringChannel{Err: nil, Str: tmp}
}

func Delpfp(id int64, sch chan StringChannel){
	db,err := openDb()
	if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	res,err := db.Query("SELECT pfp FROM users_pfps WHERE user_id=$1",id)
	if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	var tmp string
	if res.Next() {
		//fmt.Println("user has a pfp, removing...")
		err=res.Scan(&tmp) //why tho?
		if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
		_,err = db.Exec("DELETE FROM users_pfps WHERE user_id=$1;",id)
		if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	}else {
		tmp ="-1"
	}
	err = closeDb(db)
	if err!=nil {sch<-StringChannel{Err: err, Str: "-1"}}
	sch<-StringChannel{Err: nil, Str: tmp}
}