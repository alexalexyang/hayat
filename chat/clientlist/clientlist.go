package clientlist

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alexalexyang/hayat/config"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Listen(ws *websocket.Conn) {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		check(err)
	}

	listener := pq.NewListener(config.DBconfig, 10*time.Second, time.Minute, reportProblem)
	err := listener.Listen("events")
	check(err)

	fmt.Println("Start monitoring PostgreSQL...")
	for {
		waitForNotification(listener, ws)
	}
}

func waitForNotification(l *pq.Listener, ws *websocket.Conn) {
	counter := 0
	for {
		select {
		case n := <-l.Notify:
			// fmt.Println("Received data from channel [", n.Channel, "] :")
			data := notBeingServed{}
			_ = json.Unmarshal([]byte(n.Extra), &data)
			payload := []notBeingServed{data}
			ws.WriteJSON(payload)
		case <-time.After(120 * time.Second):
			fmt.Println("Received no events for 120 seconds, checking connection")
			fmt.Println("Counter:", counter)
			counter += 1
			go func() {
				l.Ping()
			}()
		}
	}
}

func ClientListHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("views/base.gohtml", "views/clientlist.gohtml")
	check(err)
	// t.ExecuteTemplate(w, "base", nil)
	if r.Method != http.MethodPost {
		t.ExecuteTemplate(w, "base", nil)
		return
	}

	roomid := r.FormValue("roomid")
	cookie1 := http.Cookie{
		Name:  "clientroom", // Set to roomid
		Value: roomid,
		// Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		// Secure:   true,
		// MaxAge:   50000,
		Path: "/",
	}
	http.SetCookie(w, &cookie1)

	cookie2 := http.Cookie{
		Name:  "consultant",
		Value: "consultant",
		// Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		// Secure:   true,
		// MaxAge:   50000,
		Path: "/",
	}
	http.SetCookie(w, &cookie2)
}

type notBeingServed struct {
	Roomid      string `json:"roomid"`
	Beingserved bool   `json:"beingserved"`
}

func ClientListWSHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()

	// Find by organisation too.
	statement := `SELECT roomid FROM rooms WHERE beingserved = 'f';`
	rows, err := db.Query(statement)
	check(err)
	defer rows.Close()

	var serviceSlice []notBeingServed

	for rows.Next() {
		serveme := notBeingServed{}
		err = rows.Scan(&serveme.Roomid)
		check(err)
		serveme.Beingserved = false
		serviceSlice = append(serviceSlice, serveme)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)
	defer ws.Close()

	ws.WriteJSON(serviceSlice)

	Listen(ws)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/base.gohtml", "views/dashboard.gohtml")
	check(err)
	// t.ExecuteTemplate(w, "base", nil)
	if r.Method != http.MethodPost {
		t.ExecuteTemplate(w, "base", nil)
		return
	}
}
