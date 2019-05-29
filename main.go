package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"database/sql"

	"github.com/alexalexyang/hayat/chat"
	"github.com/alexalexyang/hayat/config"
	"github.com/alexalexyang/hayat/login"
	"github.com/alexalexyang/hayat/models"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func main() {
	models.DBSetup()
	go listen()
	log.Println("http server started on", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, initRouter()))
}

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/anteroom", chat.AnteroomHandler)
	router.HandleFunc("/chatclientws", chat.ChatClientWSHandler)
	router.HandleFunc("/chatclient/{id:[\\w\\-]+}", chat.ChatClientHandler)
	router.HandleFunc("/clientlist", chat.ClientListHandler)

	router.HandleFunc("/login", login.loginHandler)
	return router
}

func listen() {
	_, err := sql.Open(config.Driver, config.DBconfig)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(config.DBconfig, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		panic(err)
	}

	fmt.Println("Start monitoring PostgreSQL...")
	for {
		waitForNotification(listener)
	}
}

func waitForNotification(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			fmt.Println("Received data from channel [", n.Channel, "] :")
			// Prepare notification payload for pretty print
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
			if err != nil {
				fmt.Println("Error processing JSON: ", err)
				return
			}
			fmt.Println(string(prettyJSON.Bytes()))
			return
		case <-time.After(90 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}
