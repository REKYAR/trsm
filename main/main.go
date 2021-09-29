package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)
const staticTemplate= "static/templates/"

func main() {
	cdir , err := os.Getwd()
	r := mux.NewRouter()
	r.HandleFunc("/", feedHandler)
	r.HandleFunc("/register", RegisterHandler)
	r.HandleFunc("/logout", LogoutHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/delete_account", DeleteAccHandler)
	r.HandleFunc("/profile/{id}", ProfileHandler)
	r.HandleFunc("/edit_profile", ProfileEditHandler)
	r.HandleFunc("/change_password", PasswordEditHandler)
	r.HandleFunc("/change_pfp", PfpEditHandler)
	r.PathPrefix("/pics/").Handler(http.StripPrefix("/pics/",http.FileServer(http.Dir(cdir)) ))
	err = http.ListenAndServe("localhost:8080", r)
	log.Fatal(err)
}
