package chat

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/alexalexyang/hayat/serverconfig"
	"github.com/gofrs/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var ChatroomRegistry = make(map[string]ChatroomStruct)

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

	unencodedValue := uuid.Must(uuid.NewV4()).String()
	fmt.Println("CookieValue1:", unencodedValue)
	encodedValue := encoder(unencodedValue)
	cookieSetter(w, "hayatclient", encodedValue)

	anteroomValues := AnteroomStruct{}
	anteroomValues.username = r.FormValue("username")
	anteroomValues.age, _ = strconv.Atoi(r.FormValue("age"))
	anteroomValues.gender = r.FormValue("gender")
	anteroomValues.issues = r.FormValue("issues")
	anteroomValues.token = r.FormValue("token")
	fmt.Println(anteroomValues)

	// chatroomID := uuid.Must(uuid.NewV4()).String()
	unencodedRoom := uuid.Must(uuid.NewV4()).String()
	fmt.Println("RoomValue1:", unencodedRoom)
	encodedRoom := encoder(unencodedRoom)
	cookieSetter(w, "clientroom", encodedRoom)

	// Add anteroomValues to db.

	// http.Redirect(w, r, "http://localhost:3000/contact", http.StatusFound)
	http.Redirect(w, r, serverconfig.Domain+serverconfig.Port+"/chatclient/"+encodedRoom, http.StatusSeeOther)
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
	ChatroomRegistry[chatroomID] = Chatroom
	return clients
}

func ChatClientWSHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("hayatclient")
	encodedValue := cookie.Value
	check(err)
	cookieValue := decoder(encodedValue)
	fmt.Println("CookieValue2:", cookieValue)

	roomCookie, err := r.Cookie("clientroom")
	encodedRoom := roomCookie.Value
	check(err)
	roomValue := decoder(encodedRoom)
	fmt.Println("RoomValue2:", roomValue)

	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)
	defer ws.Close()

	// Create a chatroom.
	clients := ChatRoomMaker(roomValue, ws)

	for k, v := range ChatroomRegistry {
		fmt.Println(k, v)
	}

	// Use cookie to find the user and enter into chatBroker.

	// Has to become chatBroker(clients, ws, InstanceUser) so we can send out his username.
	chatBroker(clients, ws)
}

func chatBroker(clients map[*websocket.Conn]bool, ws *websocket.Conn) {
	var msg Message
	for {
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

func ClientListHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/base.gohtml", "views/clientlist.gohtml")
	check(err)

	t.ExecuteTemplate(w, "base", nil)
}
