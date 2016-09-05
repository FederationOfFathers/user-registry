package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db = &database{}

func init() {
}

func mindSQL() {
	conn, err := sql.Open("mysql", sqlURI)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	db.DB = conn
	db.prepare()
}
