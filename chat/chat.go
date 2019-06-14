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

var RoomsRegistry = Registry{
	Rooms: make(map[string]ChatroomStruct),
}

// var RoomsRegistry.Rooms = make(map[string]ChatroomStruct)

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
	t, err := template.ParseFiles("views/base.gohtml", "views/chatclient.gohtml")
	check(err)

	t.ExecuteTemplate(w, "base", nil)
}

func ChatRoomMaker(chatroomID string, ws *websocket.Conn) ChatroomStruct {
	var Chatroom ChatroomStruct
	Chatroom.Clients = make(map[*websocket.Conn]bool)
	Chatroom.Clients[ws] = true
	RoomsRegistry.Rooms[chatroomID] = Chatroom
	return Chatroom
}

func ChatClientWSHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	urlid := params["id"]

	cookies := r.Cookies()
	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}

	if _, ok := cookieMap["consultant"]; ok {
		if _, ok := RoomsRegistry.Rooms[urlid]; ok {
			room := RoomsRegistry.Rooms[urlid]

			ws, err := upgrader.Upgrade(w, r, nil)
			check(err)
			defer ws.Close()
			room.Clients[ws] = true

			chatBroker(room, ws, cookieMap["consultantName"])
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
	if room, ok := RoomsRegistry.Rooms[roomCookie.Value]; ok {
		fmt.Println("Reconnecting.")
		room.Clients[ws] = true
		chatBroker(room, ws, username)
		return
	}

	// Create a chatroom.
	chatroom := ChatRoomMaker(roomCookie.Value, ws)

	chatBroker(chatroom, ws, username)
}

func chatBroker(room ChatroomStruct, ws *websocket.Conn, username string) {
	var msg Message
	msg.Username = username
	for {
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(room.Clients, ws)

			// Set EmptySince to time.Now() so that cleanUpRooms() can delete the room in a couple of hours.
			// if len(room.Clients) == 0 {
			// 	fmt.Println("Next line is room Emptyslice")
			// 	fmt.Println(RoomsRegistry.Rooms[room.ID].EmptySince)
			// 	delete(RoomsRegistry.Rooms, room.ID)
			// }
			break
		}
		// Send the newly received message to the broadcast channel
		for client := range room.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(room.Clients, ws)
				// if len(room.Clients) == 0 {
				// 	delete(RoomsRegistry.Rooms, room.ID)
				// 	fmt.Println("Next line is room Emptyslice")
				// 	fmt.Println(RoomsRegistry.Rooms[room.ID].EmptySince)
				// }
				break
			}
		}
	}
}

// In case program crashes, rebuild RoomsRegistry.Rooms on restart.
func Rebuild() {
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
			ID: roomid,
		}
		RoomsRegistry.Rooms[roomid] = room
	}
}

// Deletes rooms from the room registry if they've been empty for an hour.
func (r *Registry) CleanUpRooms() {
	for {
		// for room, _ := range r.Rooms {
		// 	if len(r.Rooms[room].Clients) == 0 {
		// 		fmt.Println("Next line is in cleanup")
		// 		fmt.Println(r.Rooms[room].EmptySince)
		// 		fmt.Println("Room is: ", room)

		// 		if r.Rooms[room].EmptySince.String() == "0001-01-01 00:00:00 +0000 UTC" {
		// 			fmt.Println("Wrong time.")
		// 			continue
		// 		}
		// 		if time.Since(r.Rooms[room].EmptySince) > 10.00 {
		// 			fmt.Println("Right time.")
		// 			// fmt.Println("Empty room is empty.")
		// 			delete(r.Rooms, room)

		// 			// db, err := sql.Open(config.DBType, config.DBconfig)
		// 			// check(err)
		// 			// defer db.Close()

		// 			// statement := `DELETE FROM rooms WHERE roomid=$1;`
		// 			// _, err = db.Exec(statement, room)
		// 			// check(err)
		// 		}
		// 	}
		// }
		// for k, v := range r.Rooms {
		// 	fmt.Println(k, v)
		// }
		for k, v := range RoomsRegistry.Rooms {
			fmt.Println(k)
			fmt.Println(v.EmptySince)
			fmt.Println("\n")
		}
		time.Sleep(2 * time.Second)
		// fmt.Println("\n")
	}
}
