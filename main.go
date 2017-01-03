package main

import (
	"github.com/apokalyptik/cfg"
	"github.com/hashicorp/consul/api"
)

var sqlURI = "user:password@tcp(127.0.0.1:3306)/fofgaming?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci"
var listenOn = "0.0.0.0:8875"

func init() {
	sql := cfg.New("db")
	sql.StringVar(&sqlURI, "uri", sqlURI, "MySQL Connection URI")

	api := cfg.New("api")
	api.StringVar(&listenOn, "listen", listenOn, "Listen for api connections on")
}

var consul *api.Client

func mustConsul() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	consul = client
}

func main() {
	cfg.Parse()
	mustConsul()
	mindSQL()
	go mindSeenList()
	go mindPrivateUserList()
	mindHTTP()
}
