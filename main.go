package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IssuesForm struct {
	Issues string `json:"issues"`
	Name   string `json:"name"`
	Age    string `json:"age"`
	Gender string `json:"gender"`
}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		fmt.Println(msg)
		// Send the newly received message to the broadcast channel
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}

}

func anteroomHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	fmt.Println(string(body))
	// Set cookie
	http.Redirect(w, r, "http://google.com", http.StatusSeeOther)
}

func main() {
	log.Println("http server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", initRouter()))
}

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/anteroom", anteroomHandler).Methods("POST")
	router.HandleFunc("/chat", chatHandler)
	return router
}
