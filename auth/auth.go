package auth

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/alexalexyang/explicitauth/auth"
)

var Cfg = auth.Config{
	// Host:              "localhost",
	// Protocol:          "http://",
	// Port:              ":8000",
	DBUser:            "postgres",
	DBType:            "sqlite3",
	DBName:            "auth.db",
	DBHost:            "localhost",
	DBPort:            "5431",
	DBPw:              "1234",
	EmailID:           "your_email_address_here",
	EmailPw:           os.Getenv("your_email_password"),
	SmtpHost:          "smtp.gmail.com",
	SmtpPort:          "587",
	SessionCookieName: "SessionCookie",
}

func displayTemplate(w http.ResponseWriter, r *http.Request, baseTemplate string, pageTemplate string) {
	t, err := template.ParseFiles(baseTemplate, pageTemplate)
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, "base", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, "views/base.gohtml", "views/auth/register.gohtml")
		return
	}
	Cfg.Register(r)
}
