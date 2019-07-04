package config

import (
	"fmt"
	"os"
)

var (
	// Domain config
	Domain   = "hayat.notathoughtexperiment.me"
	Port     = ":8000"
	Protocol = "https://"

	// Db config
	DBType     = "postgres"
	DBHost     = "localhost"
	DBPort     = "5431"
	DBUser     = "postgres"
	DBPassword = "1234"
	DBName     = "hayatdb"
	DBconfig   = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)

	// Email config
	EmailID  = os.Getenv("MYEMAIL")
	EmailPw  = os.Getenv("MYPW")
	SmtpHost = "smtp.gmail.com"
	SmtpPort = "587"
)
