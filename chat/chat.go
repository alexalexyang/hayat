package chat

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/alexalexyang/hayat/config"
	"github.com/gofrs/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Rooms mapped by room ID to user websockets.
var roomsRegistry = make(map[string]ChatroomStruct)

type AnteroomStruct struct {
	username string
	age      int
	gender   string
	issues   string
	token    string
}

type ChatroomStruct struct {
	ID      string
	Clients map[*websocket.Conn]bool
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var cookieKey = "courtyard"
var hashKey = []byte("very-secret")
var s = securecookie.New(hashKey, nil)

func encoder(unencodedValue string) string {
	encoded, err := s.Encode(cookieKey, unencodedValue)
	check(err)
	return encoded
}

func decoder(encodedValue string) string {
	var cookieValue string
	s.Decode(cookieKey, encodedValue, &cookieValue)
	return cookieValue
}

func cookieSetter(w http.ResponseWriter, cookieName string, encodedValue string, roomPath string) {
	cookie := http.Cookie{
		Name:  cookieName,
		Value: encodedValue,
		// Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		// Secure:   true,
		// MaxAge:   50000,
		Path: roomPath,
	}
	http.SetCookie(w, &cookie)
}

func AnteroomHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("views/base.gohtml", "views/anteroom.gohtml")
	check(err)

	if r.Method != http.MethodPost {
		t.ExecuteTemplate(w, "base", nil)
		return
	}

	username := r.FormValue("username")
	age, _ := strconv.Atoi(r.FormValue("age"))
	gender := r.FormValue("gender")
	issues := r.FormValue("issues")

	// Token is used to identify which customer the client goes to.
	token := r.FormValue("token")

	roomID := uuid.Must(uuid.NewV4()).String()
	cookieSetter(w, "clientroom", roomID, "/chatclientws/"+roomID)

	// Add anteroomValues to db.
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()

	statement := `INSERT INTO rooms (roomid, organisation, username, age, gender, issues)
				VALUES ($1, $2, $3, $4, $5, $6);`
	_, err = db.Exec(statement, roomID, token, username, age, gender, issues)
	check(err)

	http.Redirect(w, r, config.Protocol+config.Domain+config.Port+"/chatclient/"+roomID, http.StatusSeeOther)
}

func ChatClientHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/base.gohtml", "views/chatclient.gohtml")
	check(err)

	t.ExecuteTemplate(w, "base", nil)
}

func ChatRoomMaker(chatroomID string, ws *websocket.Conn) map[*websocket.Conn]bool {
	var Chatroom ChatroomStruct
	Chatroom.Clients = make(map[*websocket.Conn]bool)
	Chatroom.Clients[ws] = true
	clients := Chatroom.Clients
	roomsRegistry[chatroomID] = Chatroom
	return clients
}

func ChatClientWSHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	urlid := params["id"]

	cookies := r.Cookies()
	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
		fmt.Println(cookieMap[cookie.Name])
	}

	// var roomIDString string
	// for k, _ := range cookieMap {
	// 	if _, ok := roomsRegistry[k]; ok {
	// 		if len(roomsRegistry[k].Clients) == 1 {
	// 			roomIDString = k
	// 		}
	// 	}
	// }

	if _, ok := cookieMap["consultant"]; ok {
		if _, ok := roomsRegistry[urlid]; ok {
			room := roomsRegistry[urlid]

			ws, err := upgrader.Upgrade(w, r, nil)
			check(err)
			defer ws.Close()
			room.Clients[ws] = true

			chatBroker(room.Clients, ws, cookieMap["consultantName"])
		}
	}

	roomCookie, err := r.Cookie("clientroom")

	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)
	defer ws.Close()

	// Create a chatroom.
	clients := ChatRoomMaker(roomCookie.Value, ws)

	// Use cookie to update client room details.
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()

	statement := `UPDATE rooms SET beingserved = $1 WHERE roomid = $2;`
	_, err = db.Exec(statement, false, roomCookie.Value)
	check(err)

	// Use cookie to get client's name from their profile.
	statement = `SELECT username from rooms WHERE roomid = $1;`
	row := db.QueryRow(statement, roomCookie.Value)
	var username string
	row.Scan(&username)
	// check(err)

	chatBroker(clients, ws, username)
}

func chatBroker(clients map[*websocket.Conn]bool, ws *websocket.Conn, username string) {
	var msg Message
	msg.Username = username
	for {
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
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

func ChatConsultantWSHandler(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)

	for _, ck := range r.Cookies() {
		if ck.Name == "hayatclient" {

		}
	}

}
