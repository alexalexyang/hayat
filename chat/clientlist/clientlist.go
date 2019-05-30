package clientlist

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alexalexyang/hayat/config"
	"github.com/lib/pq"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
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

func Listen() {

	reportProblem := func(ev pq.ListenerEventType, err error) {
		check(err)
	}

	listener := pq.NewListener(config.DBconfig, 10*time.Second, time.Minute, reportProblem)
	err := listener.Listen("events")
	check(err)

	fmt.Println("Start monitoring PostgreSQL...")
	for {
		waitForNotification(listener)
	}
}

func waitForNotification(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			fmt.Println("Received data from channel [", n.Channel, "] :")
			fmt.Println(n.Extra)
			return
		case <-time.After(120 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}
