package main

import (
	"log"
	"net/http"

	"github.com/alexalexyang/hayat/auth"

	"github.com/alexalexyang/hayat/chat"
	"github.com/alexalexyang/hayat/chat/clientlist"
	"github.com/alexalexyang/hayat/config"
	"github.com/alexalexyang/hayat/models"
	"github.com/gorilla/mux"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	models.DBSetup()
	// go clientlist.Listen()
	log.Println("http server started on", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, initRouter()))
}

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/anteroom", chat.AnteroomHandler)

	// Chat client
	router.HandleFunc("/chatclientws", chat.ChatClientWSHandler)
	router.HandleFunc("/chatclient/{id:[\\w\\-]+}", chat.ChatClientHandler)

	router.HandleFunc("/dashboard", clientlist.DashboardHandler)
	// Clientlist
	router.HandleFunc("/clientlistws", clientlist.ClientListWSHandler)
	router.HandleFunc("/clientlist", clientlist.ClientListHandler)

	// Auth
	router.HandleFunc("/register/org", auth.RegisterOrgHandler)
	router.HandleFunc("/register/user/{id:[\\w\\-]+}", auth.RegisterUserHandler)
	router.HandleFunc("/invite", auth.InviteHandler)
	router.HandleFunc("/login", auth.LoginHandler)
	router.HandleFunc("/logout", auth.LogoutHandler)
	router.HandleFunc("/update", auth.UpdateHandler)
	router.HandleFunc("/forgotpw", auth.ForgotPwHandler)
	router.HandleFunc("/changepw/{id:[\\w\\-]+}", auth.ChangePwHandler)
	router.HandleFunc("/deleteaccount", auth.DeleteAccountHandler)

	return router
}
