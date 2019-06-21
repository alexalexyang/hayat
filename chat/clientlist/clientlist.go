package clientlist

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
	"net/http/httputil"
	"github.com/alexalexyang/hayat/config"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type notBeingServed struct {
	Roomid       string `json:"roomid"`
	Beingserved  bool   `json:"beingserved"`
	Organisation string `json:"organisation"`
	Username     string `json:"username"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ClientListHandler(w http.ResponseWriter, r *http.Request) {
	var username string
	var organisation string
	if r.Method != http.MethodPost {
		t, err := template.ParseFiles("views/base.gohtml", "views/clientlist.gohtml")
		check(err)

		sessionCookie, err := r.Cookie("SessionCookie")
		if err != nil {
			check(err)
			return
		}

		db, err := sql.Open(config.DBType, config.DBconfig)
		check(err)
		defer db.Close()

		statement := `SELECT email FROM sessions WHERE cookie=$1;`
		row := db.QueryRow(statement, sessionCookie.Value)
		var email string
		switch err := row.Scan(&email); err {
		case sql.ErrNoRows:
			fmt.Println("No row found.")
		case nil:
			// fmt.Println(email)
		default:
			check(err)
		}

		statement = `SELECT username, organisation FROM users WHERE email=$1;`
		row = db.QueryRow(statement, email)
		row.Scan(&username, &organisation)

		consultantName := http.Cookie{
			Name:  "consultantName",
			Value: username,
			// Expires:  time.Now().Add(time.Hour),
			HttpOnly: true,
			// Secure:   true,
			// MaxAge:   50000,
			Path: "/",
		}
		http.SetCookie(w, &consultantName)

		consultantCookie := http.Cookie{
			Name:  "consultant",
			Value: organisation,
			// Expires:  time.Now().Add(time.Hour),
			HttpOnly: true,
			// Secure:   true,
			// MaxAge:   50000,
			Path: "/",
		}
		http.SetCookie(w, &consultantCookie)
		t.ExecuteTemplate(w, "base", nil)
		return
	}

	roomid := r.FormValue("roomid")
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()
	statement := `UPDATE rooms SET beingserved = $1 WHERE roomid = $2;`
	_, err = db.Exec(statement, true, roomid)
	check(err)
}

func ClientListWSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Clientlist WS: I got a request.")
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	consultantCookie, err := r.Cookie("consultant")
	if err != nil {
		check(err)
		return
	}
	organisation := consultantCookie.Value
	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()

	// Find by organisation too.
	statement := `SELECT roomid, username FROM rooms WHERE beingserved='f' AND organisation=$1;`
	rows, err := db.Query(statement, organisation)
	check(err)
	defer rows.Close()

	var serviceSlice []notBeingServed

	for rows.Next() {
		serveme := notBeingServed{}
		err = rows.Scan(&serveme.Roomid, &serveme.Username)
		check(err)
		serveme.Beingserved = false
		serviceSlice = append(serviceSlice, serveme)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	check(err)
	defer ws.Close()

	ws.WriteJSON(serviceSlice)

	Listen(ws, organisation)
}

func ClientProfileHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("consultant")
	if err != nil {
		check(err)
		return
	}
	t, err := template.ParseFiles("views/base.gohtml", "views/clientprofile.gohtml")
	check(err)

	params := mux.Vars(r)
	urlid := params["id"]

	db, err := sql.Open(config.DBType, config.DBconfig)
	check(err)
	defer db.Close()

	clientProfile := struct {
		Username     string
		Age          string
		Gender       string
		Issues       string
		Roomid       string
		Organisation string
	}{}

	statement := `SELECT username, age, gender, issues, roomid, organisation FROM rooms WHERE roomid=$1;`
	row := db.QueryRow(statement, urlid)
	row.Scan(&clientProfile.Username, &clientProfile.Age, &clientProfile.Gender, &clientProfile.Issues, &clientProfile.Roomid, &clientProfile.Organisation)

	t.ExecuteTemplate(w, "base", clientProfile)
}

func Listen(ws *websocket.Conn, organisation string) {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		check(err)
	}

	listener := pq.NewListener(config.DBconfig, 10*time.Second, time.Minute, reportProblem)
	err := listener.Listen("events")
	check(err)

	fmt.Println("Start monitoring PostgreSQL...")
	for {
		waitForNotification(listener, ws, organisation)
	}
}

func waitForNotification(l *pq.Listener, ws *websocket.Conn, organisation string) {
	counter := 0
	for {
		select {
		case n := <-l.Notify:
			data := notBeingServed{}

			_ = json.Unmarshal([]byte(n.Extra), &data)
			fmt.Println("n.Extra:")
			fmt.Println(n.Extra)
			if organisation != data.Organisation {
				return
			}
			payload := []notBeingServed{data}
			ws.WriteJSON(payload)
		case <-time.After(120 * time.Second):
			fmt.Println("Received no events for 120 seconds, checking connection")
			l.Ping()
			fmt.Println(counter)
			counter++
			// go func() {
			// 	l.Ping()
			// }()
		}
	}
}
