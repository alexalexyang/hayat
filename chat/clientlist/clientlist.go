package clientlist

import (
	"database/sql"
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

	ws.WriteJSON("testties")
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
	for {
		select {
		case n := <-l.Notify:
			fmt.Println("Received data from channel [", n.Channel, "] :")
			fmt.Println(n.Extra)
			ws.WriteJSON(n.Extra)
		case <-time.After(120 * time.Second):
			fmt.Println("Received no events for 120 seconds, checking connection")
			go func() {
				l.Ping()
			}()
		}
	}
}

func ClientListHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/base.gohtml", "views/clientlist.gohtml")
	check(err)

	// Match session cookie in CUSTOMERS TABLE to find out which organisation consultant is from.

	// Using organisation, look in ROOMS TABLE to find all available chatrooms.
	// If beingserved=true, display. Else, don't.

	// When consultant clicks on roomID, redirect to chatclient/roomID.
	// Set roomID as cookie for that specific room.
	// And connect to another consultant's version of chatClientWSHandler.
	// In it, get roomID cookie, use it to find the room in roomsRegistry.
	// Add consultant's websocket to that room.
	// Launch chatBroker.

	t.ExecuteTemplate(w, "base", nil)
}

func ClientListWSHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open(config.Driver, config.DBconfig)
	check(err)
	defer db.Close()

	// Find by organisation too.
	statement := `SELECT roomid FROM rooms WHERE beingserved = 'f';`
	rows, err := db.Query(statement)
	check(err)
	defer rows.Close()

	notBeingServed := make(map[string]bool)

	for rows.Next() {
		var roomid string
		err = rows.Scan(&roomid)
		check(err)
		// beingserved := strconv.FormatBool(beingservedBool)

		notBeingServed[roomid] = false
	}

	for k, v := range notBeingServed {
		fmt.Println(k[:8], v)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)
	defer ws.Close()

	ws.WriteJSON(notBeingServed)

	Listen(ws)
}
