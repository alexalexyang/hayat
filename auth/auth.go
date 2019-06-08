package auth

import (
	"log"
	"net/http"
	"text/template"

	"github.com/alexalexyang/hayat/config"

	"github.com/alexalexyang/explicitauth/auth"
)

var Cfg = auth.Config{
	Host:              config.Domain,
	Protocol:          config.Protocol,
	Port:              config.Port,
	DBUser:            config.DBUser,
	DBType:            config.DBType,
	DBName:            config.DBName,
	DBHost:            config.DBHost,
	DBPort:            config.DBPort,
	DBPw:              config.DBPassword,
	EmailID:           config.EmailID,
	EmailPw:           config.EmailPw,
	SmtpHost:          config.SmtpHost,
	SmtpPort:          config.SmtpPort,
	SessionCookieName: "SessionCookie",
}

func displayTemplate(w http.ResponseWriter, r *http.Request, baseTemplate string, pageTemplate string) {
	t, err := template.ParseFiles(baseTemplate, pageTemplate)
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, "base", nil)
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, "views/base.gohtml", "views/auth/deleteaccount.gohtml")
		return
	}
	Cfg.DeleteAccount(w, r)
}

func ChangePwHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, "views/base.gohtml", "views/auth/changepw.gohtml")
		return
	}
	Cfg.ChangePw(w, r)
}

func ForgotPwHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, "views/base.gohtml", "views/auth/forgotpw.gohtml")
		return
	}
	Cfg.ForgotPw(w, r)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, "views/base.gohtml", "views/auth/logout.gohtml")
		return
	}
	Cfg.Logout(w, r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, "views/base.gohtml", "views/auth/login.gohtml")
		return
	}
	Cfg.Login(w, r)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, "views/base.gohtml", "views/auth/register.gohtml")
		return
	}
	Cfg.Register(r)
}
