package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexalexyang/hayat/auth"

	"github.com/alexalexyang/hayat/chat"
	"github.com/alexalexyang/hayat/chat/dashboard"
	"github.com/alexalexyang/hayat/generic"
	"github.com/alexalexyang/hayat/config"
	"github.com/alexalexyang/hayat/models"
	"github.com/gorilla/mux"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	models.DBSetup()

	RoomsRegistry := chat.Registry{
		Rooms: make(map[string]*chat.ChatroomStruct),
	}

	RoomsRegistry.Rebuild()
	go RoomsRegistry.CleanUpRooms()

	log.Println("http server started on", config.Port)
	http.ListenAndServe(config.Port, initRouter(&RoomsRegistry))
}

func initRouter(rg *chat.Registry) *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/anteroom", chat.AnteroomHandler)

	// Chat client
	router.HandleFunc("/chatclientws/{id:[\\w\\-]+}", rg.ChatClientWSHandler)
	router.HandleFunc("/chatclient/{id:[\\w\\-]+}", chat.ChatClientHandler)

	// Dashboard
	router.HandleFunc("/dashboardws", dashboard.DashboardWSHandler)
	router.HandleFunc("/dashboard", dashboard.DashboardHandler)
	router.HandleFunc("/clientprofile/{id:[\\w\\-]+}", dashboard.ClientProfileHandler)

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

	// Hayat pages
	router.HandleFunc("/", generic.MainPageHandler)

	return router
}
