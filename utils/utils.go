package utils

import (
	"log"
	"math/rand"
	"net/http"
)

func ErrCheck(err error)  {
	if err!=nil {
		log.Fatal(err)
	}
}

func Find(slice *[]string, val string) (int, bool) {
	for i, item := range *slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}


func CheckForCookie(request *http.Request, name string ) bool {
	_,err:=request.Cookie(name)
	if err!=nil {
		return false
	}
	return true
}

func RemoveFromSlice(slc *[]string, val string){
	idx,is := Find(slc,val)
	if is {
		*slc = append((*slc)[:idx], (*slc)[idx+1:]...)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")


func GenerateRandomName(n int) string{
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//func executeHtml(html *template.Template, writer http.ResponseWriter,data interface{}){
//	if sessions {
//
//	}
//}