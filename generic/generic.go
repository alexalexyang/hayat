package generic

import (
	"fmt"
	"net/http"
	"html/template"
)

type LoggedInDetails struct {
	Username string `json:"username"`
	LoggedIn bool `json:"isadmin"`
	Role string `json:"role"`
	Organisation string `json:"organisation"`
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func SetCookieHayat(w http.ResponseWriter, cookieName string, encodedValue string, roomPath string) {
	cookie := http.Cookie{
		Name:  cookieName,
		Value: encodedValue,
		// Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   true,
		// MaxAge:   50000,
		Path: roomPath,
	}
	http.SetCookie(w, &cookie)
}

func GetAllCookies(r *http.Request) map[string]string {
	cookies := r.Cookies()
	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}
	return cookieMap
}

func MakeLoggedInPayload(r *http.Request) LoggedInDetails {
	cookieMap := GetAllCookies(r)

	var username string
	var organisation string
	var role string

	if _, ok := cookieMap["SessionCookie"]; ok {
		username = cookieMap["username"]
		organisation = cookieMap["organisation"]
		role = cookieMap["role"]
		loggedInPayload := LoggedInDetails {
			LoggedIn: true,
			Username: username,
			Organisation: organisation,
			Role: role,
		}
	
		return loggedInPayload
	}

	return LoggedInDetails{}
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/base.gohtml", "views/navbar.gohtml", "views/generic/landing.gohtml")
	check(err)

	t.ExecuteTemplate(w, "base", nil)
	return
}