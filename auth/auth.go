package auth

import (
	"log"
	"net/http"
	"text/template"

	"github.com/alexalexyang/hayat/config"
	"github.com/alexalexyang/hayat/generic"

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

func displayTemplate(w http.ResponseWriter, r *http.Request, data generic.LoggedInDetails, baseTemplate string, navbarTemplate string, pageTemplate string) {
	t, err := template.ParseFiles(baseTemplate, navbarTemplate, pageTemplate)
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, "base", data)
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, generic.MakeLoggedInPayload(r), "./views/base.gohtml", "./views/navbar.gohtml", "./views/auth/deleteaccount.gohtml")
		return
	}
	if Cfg.DeleteAccount(w, r) == false {
		w.WriteHeader(http.StatusResetContent)
		return
	}
	
	http.Redirect(w, r, config.Protocol+config.Domain+"/login", http.StatusSeeOther)
}

func ChangePwHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, generic.MakeLoggedInPayload(r), "./views/base.gohtml", "./views/navbar.gohtml", "./views/auth/changepw.gohtml")
		return
	}
	if Cfg.ChangePw(w, r) == false {
		w.WriteHeader(http.StatusResetContent)
		return
	}
	
	http.Redirect(w, r, config.Protocol+config.Domain+"/login", http.StatusSeeOther)
}

func ForgotPwHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, generic.MakeLoggedInPayload(r), "./views/base.gohtml", "./views/navbar.gohtml", "./views/auth/forgotpw.gohtml")
		return
	}
	if Cfg.ForgotPw(w, r) == false {
		w.WriteHeader(http.StatusResetContent)
		return
	}
	
	http.Redirect(w, r, config.Protocol+config.Domain+"/login", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, generic.MakeLoggedInPayload(r), "./views/base.gohtml", "./views/navbar.gohtml", "./views/auth/logout.gohtml")
		return
	}
	if Cfg.Logout(w, r) == false {
		w.WriteHeader(http.StatusResetContent)
		return
	}
	
	http.Redirect(w, r, config.Protocol+config.Domain+"/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, generic.MakeLoggedInPayload(r), "./views/base.gohtml", "./views/navbar.gohtml", "./views/auth/login.gohtml")
		return
	}
	if Cfg.Login(w, r) == false {
		w.WriteHeader(http.StatusResetContent)
		return
	}

	
	http.Redirect(w, r, config.Protocol+config.Domain+"/dashboard", http.StatusSeeOther)
}

func RegisterOrgHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, generic.MakeLoggedInPayload(r), "./views/base.gohtml", "./views/navbar.gohtml", "./views/auth/registerorg.gohtml")
		return
	}
	if Cfg.Register(r) == false {
		w.WriteHeader(http.StatusResetContent)
		return
	}
	
	http.Redirect(w, r, config.Protocol+config.Domain+"/login", http.StatusSeeOther)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, generic.MakeLoggedInPayload(r), "./views/base.gohtml", "./views/navbar.gohtml", "./views/auth/registeruser.gohtml")
		return
	}
	if Cfg.Register(r) == false {
		w.WriteHeader(http.StatusResetContent)
		return
	}
	
	http.Redirect(w, r, config.Protocol+config.Domain+"/login", http.StatusSeeOther)
}

func InviteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, generic.MakeLoggedInPayload(r), "./views/base.gohtml", "./views/navbar.gohtml", "./views/auth/invite.gohtml")
		return
	}
	if Cfg.Invite(w, r) == false {
		w.WriteHeader(http.StatusResetContent)
		return
	}
	
	http.Redirect(w, r, config.Protocol+config.Domain+"/dashboard", http.StatusSeeOther)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		displayTemplate(w, r, generic.MakeLoggedInPayload(r), "./views/base.gohtml", "./views/navbar.gohtml", "./views/auth/update.gohtml")
		return
	}
	if Cfg.Update(r) == false {
		w.WriteHeader(http.StatusResetContent)
		return
	}
	
	http.Redirect(w, r, config.Protocol+config.Domain+"/dashboard", http.StatusSeeOther)
}
