package main

import (
	"github.com/FederationOfFathers/consul"
	"github.com/apokalyptik/cfg"
)

var sqlURI = "user:password@tcp(127.0.0.1:3306)/fofgaming?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci"
var listenOn = "0.0.0.0:8875"

func init() {
	sql := cfg.New("db")
	sql.StringVar(&sqlURI, "uri", sqlURI, "MySQL Connection URI")

	api := cfg.New("api")
	api.StringVar(&listenOn, "listen", listenOn, "Listen for api connections on")
}

func main() {
	cfg.Parse()
	consul.Must()
	mindSQL()
	go loadUsersFromDatabase()
	go mindDiscord()
	mindHTTP()
}
