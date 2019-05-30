package chat

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/alexalexyang/hayat/config"
	"github.com/gofrs/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
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

func cookieSetter(w http.ResponseWriter, cookieName string, encodedValue string) {
	cookie := http.Cookie{
		Name:  cookieName,
		Value: encodedValue,
		// Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		// Secure:   true,
		// MaxAge:   50000,
		Path: "/",
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

	// This is the unique ID for the client.
	rawCookieValue := uuid.Must(uuid.NewV4()).String()
	encodedValue := encoder(rawCookieValue)
	cookieSetter(w, "hayatclient", encodedValue)

	username := r.FormValue("username")
	age, _ := strconv.Atoi(r.FormValue("age"))
	gender := r.FormValue("gender")
	issues := r.FormValue("issues")

	// Token is used to identify which customer the client goes to.
	token := r.FormValue("token")

	// chatroomID is used to list the room in the clientlist.
	roomID := uuid.Must(uuid.NewV4()).String()
	cookieSetter(w, "clientroom", roomID)

	// Add anteroomValues to db.
	db, err := sql.Open(config.Driver, config.DBconfig)
	check(err)
	defer db.Close()

	statement := `INSERT INTO clientprofiles (sessioncookie, username, age, gender, issues)
	VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(statement, rawCookieValue, username, age, gender, issues)
	check(err)

	statement = `INSERT INTO rooms (roomid, organisation, sessioncookie)
				VALUES ($1, $2, $3);`
	_, err = db.Exec(statement, roomID, token, rawCookieValue)
	check(err)

	// http.Redirect(w, r, "http://localhost:3000/contact", http.StatusFound)
	http.Redirect(w, r, config.Domain+config.Port+"/chatclient/"+roomID, http.StatusSeeOther)
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
	cookie, err := r.Cookie("hayatclient")
	encodedValue := cookie.Value
	check(err)
	cookieValue := decoder(encodedValue)

	roomCookie, err := r.Cookie("clientroom")

	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)
	defer ws.Close()

	// Create a chatroom.
	clients := ChatRoomMaker(roomCookie.Value, ws)

	// Use cookie to find the user and enter into chatBroker.
	db, err := sql.Open(config.Driver, config.DBconfig)
	check(err)
	defer db.Close()

	statement := `UPDATE rooms SET beingserved = $1 WHERE sessioncookie = $2;`
	_, err = db.Exec(statement, false, cookieValue)
	check(err)

	statement = `SELECT username from clientprofiles WHERE sessioncookie = $1;`
	row := db.QueryRow(statement, cookieValue)
	var username string
	row.Scan(&username)
	// check(err)

	// Has to become chatBroker(clients, ws, InstanceUser) so we can send out his username.
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
