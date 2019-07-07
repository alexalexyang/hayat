package chat

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/alexalexyang/hayat/config"
	"github.com/alexalexyang/hayat/generic"
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
type Registry struct {
	Rooms map[string]*ChatroomStruct
}

type AnteroomStruct struct {
	username string
	age      int
	gender   string
	issues   string
	token    string
}

type ChatroomStruct struct {
	ID         string
	EmptySince time.Time
	Clients    map[*websocket.Conn]bool
}

type Message struct {
	RoomID	string `json:"roomid"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Type string `json:"type"`
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

func AnteroomHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("./views/base.gohtml", "./views/navbar.gohtml", "./views/anteroom.gohtml")
	check(err)

	if r.Method != http.MethodPost {
		t.ExecuteTemplate(w, "anteroom", nil)
		return
	}

	username := r.FormValue("username")
	age, _ := strconv.Atoi(r.FormValue("age"))
	gender := r.FormValue("gender")
	issues := r.FormValue("issues")

	// Token is used to identify which customer the client goes to.
	token := r.FormValue("token")

	roomID := uuid.Must(uuid.NewV4()).String()
	generic.SetCookieHayat(w, "clientroom", roomID, "/chatclientws/"+roomID)

	// Set client username as cookie so we can use it set up the chatBroker in the chatroom.
	generic.SetCookieHayat(w, "clientusername", username, "/chatclientws/"+roomID)

	// Add anteroomValues to db.
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()

	statement := `INSERT INTO rooms (roomid, organisation, username, age, gender, issues, beingserved)
				VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err = db.Exec(statement, roomID, token, username, age, gender, issues, false)
	check(err)

	http.Redirect(w, r, config.Protocol+config.Domain+"/chatclient/"+roomID, http.StatusSeeOther)
}

func ChatClientHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./views/chatclient.gohtml")
	check(err)

	t.ExecuteTemplate(w, "chatclient", nil)
}

func (rg *Registry) makeChatroom(chatroomID *string) ChatroomStruct {
	var Chatroom ChatroomStruct
	Chatroom.Clients = make(map[*websocket.Conn]bool)
	Chatroom.ID = *chatroomID
	return Chatroom
}

func (rg *Registry) getRoom(roomid string) *ChatroomStruct {
	return rg.Rooms[roomid]
}

func (rg *Registry) ChatClientWSHandler(w http.ResponseWriter, r *http.Request) {
	// Open db connection because we're going to use it a few times.
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()

	// Get the room ID from URL parameters.
	params := mux.Vars(r)
	roomid := params["id"]

	// Get all cookies for the room
	cookieMap := generic.GetAllCookies(r)

	// Upgrade connection so we can send chat history up if it exists.
	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)
	defer ws.Close()

	// Check database for chatroom. If exists, reload chat history.
	statement := `SELECT username, message, type FROM messages WHERE roomid=$1;`
	rows, err := db.Query(statement, roomid)
	check(err)
	defer rows.Close()

	payload := []Message{}
	for rows.Next() {
		var savedName string
		var savedMsg string
		var savedType string
		err = rows.Scan(&savedName, &savedMsg, &savedType)
		check(err)
		msg := Message{
			Username: savedName,
			Message: savedMsg,
			Type: savedType,
		}
		payload = append(payload, msg)
	}

	ws.WriteJSON(payload)

	// Reconnect or is consultant.

	if _, ok := rg.Rooms[roomid]; ok {
	
		room := rg.Rooms[roomid]
		room.Clients[ws] = true

		var username string

		if consultantUsername, ok := cookieMap["consultantName"]; ok {
			username = consultantUsername
			rg.chatBroker(room, ws, username)
			return
		}

		if clientUsername, ok := cookieMap["clientusername"]; ok {
			username = clientUsername
			rg.chatBroker(room, ws, username)
			return
		}

		// rg.chatBroker(room, ws, username)
		// return
	}

	// Make the chatroom.
	chatroom := rg.makeChatroom(&roomid)

	// Register the room.
	rg.Rooms[roomid] = &chatroom
	
	// Add client to chatroom.
	chatroom.Clients[ws] = true
	
	// Launch the chatBroker.
	rg.chatBroker(&chatroom, ws, cookieMap["clientusername"])
}

func Ping(ws *websocket.Conn) {
	for {
		if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			return
		}
		time.Sleep(50 * time.Second)
	}
}

func saveMsg(roomid string, msg Message) {
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()
	statement := `INSERT INTO messages (timestamptz, roomid, username, message, type)
	VALUES ($1, $2, $3, $4, $5);`
	_, err = db.Exec(statement, time.Now(), roomid, msg.Username, msg.Message, msg.Type)
	check(err)
}

// Update room emptysince field in database so cleanup can delete after an hour.
func setEmpty(roomid string) {
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()
	statement := `UPDATE rooms SET emptysince = $1 WHERE roomid = $2;`
	_, err = db.Exec(statement, time.Now(), roomid)
	check(err)
}

func (rg *Registry) chatBroker(room *ChatroomStruct, ws *websocket.Conn, username string) {
	var msg Message
	msg.Username = username
	msg.Type = "open"
	msg.RoomID = room.ID
	saveMsg(room.ID, msg)
	payload := []Message{msg}
	for client := range room.Clients {
		err := client.WriteJSON(&payload)
		check(err)
		msg.Type = ""
	}

	go Ping(ws)

	for {

		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %v", err)
			msg.Type = "close"
			payload = []Message{msg}
			for client := range room.Clients {
				err := client.WriteJSON(&payload)
				check(err)
			}
			saveMsg(room.ID, msg)
			ws.Close()
			delete(room.Clients, ws)
			setEmpty(room.ID)
			return
		}

		// Send the newly received message to the broadcast channel
		for client := range room.Clients {
			payload = []Message{msg}
			err := client.WriteJSON(payload)
			if err != nil {
				log.Printf("error: %v", err)
				msg.Type = "close"
				payload := []Message{msg}
				err := ws.WriteJSON(&payload)
				check(err)
				saveMsg(room.ID, msg)
				ws.Close()
				delete(room.Clients, ws)
				setEmpty(room.ID)
				return
			}
		}
		// Save message to database.
		saveMsg(room.ID, msg)
	}
}

// In case program crashes, rebuild RoomsRegistry.Rooms on restart.
func (rg *Registry) Rebuild() {
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()

	statement := `SELECT roomid FROM rooms;`
	rows, err := db.Query(statement)
	check(err)

	for rows.Next() {
		var roomid string
		err = rows.Scan(&roomid)
		check(err)

		room := ChatroomStruct{
			ID:      roomid,
			Clients: make(map[*websocket.Conn]bool),
		}
		rg.Rooms[roomid] = &room
	}

}

// Deletes rooms from the room registry if they've been empty for an hour.
func (r *Registry) CleanUpRooms() {
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()
	for {
		for room, _ := range r.Rooms {
			if len(r.Rooms[room].Clients) == 0 {
				statement := `SELECT emptysince from rooms WHERE roomid = $1;`
				row := db.QueryRow(statement, room)
				var emptysince time.Time
				row.Scan(&emptysince)
				if time.Since(emptysince).Seconds() > 3600.00 {
					delete(r.Rooms, room)
					statement := `DELETE FROM rooms WHERE roomid=$1;`
					_, err = db.Exec(statement, room)
					check(err)
					statement = `DELETE FROM messages WHERE roomid=$1;`
					_, err = db.Exec(statement, room)
					check(err)

					delete(r.Rooms, room)
				}
				continue
			}
			

		}
		time.Sleep(3600 * time.Second)
	}
}
