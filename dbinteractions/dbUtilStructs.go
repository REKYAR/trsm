package dbinteractions

import "database/sql"

type IdChannel struct {
	Id int64
	Err error
}

type StringChannel struct {
	Str string
	Err error
}

//type User struct {
//	LoggedIn bool
//	Id int
//	Login string
//	TagWhitelist string
//	TagBlacklist string
//	Rep int
//	Password string
//	Admin bool
//}

type  UserChannel struct {
	UserField *sql.Rows
	Err error
}
