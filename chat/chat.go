package chat

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/alexalexyang/hayat/serverconfig"
	"github.com/gofrs/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
)

var ChatroomsRegistry = make(map[string]ChatroomsStruct)

type ChatroomsStruct struct {
	Clients map[*websocket.Conn]bool
}

type Message struct {
	Email    string `json:"email"`
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

func AnteroomHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/base.gohtml", "views/anteroom.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	if r.Method != http.MethodPost {
		t.ExecuteTemplate(w, "base", nil)
		return
	}

	jsonMap := make(map[string]interface{})

	jsonMap["FirstName"] = r.FormValue("first_name")
	jsonMap["LastName"] = r.FormValue("last_name")
	jsonMap["Age"] = r.FormValue("age")
	jsonMap["Gender"] = r.FormValue("gender")
	jsonMap["Issues"] = r.FormValue("issues")

	fmt.Println(jsonMap)

	// http.Redirect(w, r, "http://localhost:3000/contact", http.StatusFound)
	http.Redirect(w, r, serverconfig.Domain+serverconfig.Port+"/chatclient", http.StatusSeeOther)
}

var c = make(chan string)

func ChatClientWSHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("hayat")
	if err != nil {
		log.Fatal(err)
	}
	var cookieValue string
	s.Decode("hayat", cookie.Value, &cookieValue)
	fmt.Println(cookieValue)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	var Chatrooms ChatroomsStruct
	Chatrooms.Clients = make(map[*websocket.Conn]bool)
	Chatrooms.Clients[ws] = true
	clients := Chatrooms.Clients
	ChatroomsRegistry[cookieValue] = Chatrooms

	c <- cookieValue

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

var hashKey = []byte("very-secret")
var s = securecookie.New(hashKey, nil)

func cookieHandler(w http.ResponseWriter) {
	u1 := uuid.Must(uuid.NewV4())
	fmt.Println(u1)
	encoded, err := s.Encode("hayat", u1.String())
	if err != nil {
		log.Fatal(err)
	}
	cookie := http.Cookie{
		Name:  "hayat",
		Value: encoded,
		// Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		// Secure:   true,
		// MaxAge:   50000,
		Path: "/",
	}
	http.SetCookie(w, &cookie)
}

func ChatClientHandler(w http.ResponseWriter, r *http.Request) {

	cookieHandler(w)

	t, err := template.ParseFiles("views/base.gohtml", "views/chatclient.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	t.ExecuteTemplate(w, "base", nil)

}

func ClientListWSHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	result := <-c
	fmt.Println(result)
	fmt.Println(ChatroomsRegistry)
	rtest := struct {
		Integer string `json:"int"`
	}{
		result,
	}
	ws.WriteJSON(rtest)
}

func ClientListHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/base.gohtml", "views/clientlist.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	t.ExecuteTemplate(w, "base", nil)
}
