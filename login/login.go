package login

import (
	"log"
	"net/http"
	"text/template"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/base.gohtml", "views/login/login.gohtml")
	check(err)

	// Set session cookie.

	// Write session cookie to customers db.

	t.ExecuteTemplate(w, "base", nil)
}
