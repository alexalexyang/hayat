package config

import "fmt"

// Server config
var Domain = "http://localhost"
var Port = ":8000"

// Db config
var Driver = "postgres"

const (
	host     = "localhost"
	port     = 5431
	user     = "postgres"
	password = "1234"
	dbname   = "hayatdb"
)

var DBconfig = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
