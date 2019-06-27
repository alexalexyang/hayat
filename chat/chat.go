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
	Rooms map[string]ChatroomStruct
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
	Username string `json:"username"`
	Message  string `json:"message"`
	Timestamp time.Time `json:"timestamp"`
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

func CookieSetter(w http.ResponseWriter, cookieName string, encodedValue string, roomPath string) {
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

	t, err := template.ParseFiles("./views/base.gohtml", "./views/navbar.gohtml", "./views/anteroom.gohtml")
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
	CookieSetter(w, "clientroom", roomID, "/chatclientws/"+roomID)

	// Add anteroomValues to db.
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()

	statement := `INSERT INTO rooms (roomid, organisation, username, age, gender, issues, beingserved)
				VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err = db.Exec(statement, roomID, token, username, age, gender, issues, false)
	check(err)

	http.Redirect(w, r, config.Protocol+config.Domain+config.Port+"/chatclient/"+roomID, http.StatusSeeOther)
}

func ChatClientHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./views/base.gohtml", "./views/navbar.gohtml", "./views/chatclient.gohtml")
	check(err)

	t.ExecuteTemplate(w, "base", nil)
}

func (rg *Registry) ChatRoomMaker(chatroomID *string, ws *websocket.Conn) ChatroomStruct {
	var Chatroom ChatroomStruct
	Chatroom.Clients = make(map[*websocket.Conn]bool)
	Chatroom.Clients[ws] = true
	Chatroom.ID = *chatroomID
	rg.Rooms[*chatroomID] = Chatroom
	return Chatroom
}

func (r *Registry) getRoom(roomid string) ChatroomStruct {
	return r.Rooms[roomid]
}

func (rg *Registry) ChatClientWSHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomid := params["id"]
	cookies := r.Cookies()
	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}
	
	if _, ok := cookieMap["consultant"]; ok {
		if _, ok := rg.Rooms[roomid]; ok {
			room := rg.getRoom(roomid)

			ws, err := upgrader.Upgrade(w, r, nil)
			check(err)
			defer ws.Close()
			room.Clients[ws] = true

			rg.chatBroker(&room, ws, cookieMap["consultantName"])
		}
		return
	}
	
	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)
	defer ws.Close()
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()

	roomCookie, err := r.Cookie("clientroom")
	check(err)

	// Use cookie to get client's name from their profile.
	statement := `SELECT username from rooms WHERE roomid = $1;`
	row := db.QueryRow(statement, roomCookie.Value)
	var username string
	row.Scan(&username)
	// Check if user was already in a room and just reconnecting now.
	if _, ok := rg.Rooms[roomCookie.Value]; ok {
		room := rg.getRoom(roomid)
		fmt.Println("Reconnecting.")
		room.Clients[ws] = true


		payload := []Message{}
		statement := `SELECT username, message FROM messages WHERE roomid=$1;`
		rows, err := db.Query(statement, roomid)
		check(err)
	
		for rows.Next() {
			var savedName string
			var savedMsg string
			err = rows.Scan(&savedName, &savedMsg)
			check(err)
			msg := Message{
				Username: savedName,
				Message: savedMsg,
			}
			payload = append(payload, msg)
		}
		ws.WriteJSON(payload)
		rg.chatBroker(&room, ws, username)
		return
	}
	// Create a chatroom.
	chatroom := rg.ChatRoomMaker(&roomCookie.Value, ws)
	rg.chatBroker(&chatroom, ws, username)
}

func Ping(ws *websocket.Conn) {
	for {
		if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			return
		}
		time.Sleep(50 * time.Second)
	}
}

func saveMsg(roomid string, username string, message string) {
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()
	statement := `INSERT INTO messages (timestamptz, roomid, username, message)
	VALUES ($1, $2, $3, $4);`
	_, err = db.Exec(statement, time.Now(), roomid, username, message)
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
	go Ping(ws)
	for {
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(room.Clients, ws)
			setEmpty(room.ID)
			return
		}
		// Send the newly received message to the broadcast channel
		for client := range room.Clients {
			payload := []Message{msg}
			err := client.WriteJSON(payload)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(room.Clients, ws)
				setEmpty(room.ID)
				return
			}
		}
				// Save message to database.
				saveMsg(room.ID, username, msg.Message)
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
		rg.Rooms[roomid] = room
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

					delete(r.Rooms, room)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
