package generic

import (
	"fmt"
	"net/http"
	"html/template"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/base.gohtml", "views/navbar.gohtml", "views/generic/landing.gohtml")
	check(err)

	t.ExecuteTemplate(w, "base", nil)
	return
}