package main

import (
	"log"
	"net/http"

	"github.com/alexalexyang/hayat/chat"
	"github.com/alexalexyang/hayat/serverconfig"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("http server started on", serverconfig.Port)
	log.Fatal(http.ListenAndServe(serverconfig.Port, initRouter()))
}

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/anteroom", chat.AnteroomHandler)
	router.HandleFunc("/chatclientws", chat.ChatClientWSHandler)
	router.HandleFunc("/chatclient/{id:[\\w\\-\\=]+}", chat.ChatClientHandler)
	router.HandleFunc("/clientlist", chat.ClientListHandler)
	return router
}
