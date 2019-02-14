package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var getAllUsers *sql.Stmt
var getDiscordUser *sql.Stmt
var insertUser *sql.Stmt
var updateName *sql.Stmt
var updateSeen *sql.Stmt

func mustPrepare(q string) *sql.Stmt {
	s, e := db.Prepare(q)
	if e != nil {
		panic(e)
	}
	return s
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
	db = conn
	getAllUsers = mustPrepare("SELECT `ID`,`xbl`,`destiny`,`seen`,`name`,`tz`,`discord` FROM `members` WHERE discord IS NOT NULL")
	getDiscordUser = mustPrepare("SELECT `ID`,`xbl`,`destiny`,`seen`,`name`,`tz`,`discord` FROM `members` WHERE `discord`=?")
	insertUser = mustPrepare("INSERT IGNORE INTO `members` (`discord`,`name`,`seen`,`updated_at`,`created_at`) VALUES(?,?,UNIX_TIMESTAMP(), NOW(), NOW())")
	updateName = mustPrepare("UPDATE `members` SET `name`=?, `updated_at`=NOW() WHERE `id` = ? LIMIT 1")
	updateSeen = mustPrepare("UPDATE `members` SET `seen`=?, `updated_at`=NOW() WHERE `id` = ? LIMIT 1")
}
